package engine

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"codedock.dev/codedock/internal/models"
)

func (d *Deployer) prepareServerlessCode(app *models.AppService, sourceDir string, logWriter io.Writer) error {
	if app.BuildEngine != models.BuildEngineServerless {
		return nil
	}

	code, err := d.store.GetServerlessFunctionCode(app.ID)
	if err != nil {
		return fmt.Errorf("could not retrieve serverless code: %w", err)
	}

	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		return fmt.Errorf("could not create source directory: %w", err)
	}

	if code.Runtime == "nodejs" {
		return prepareNodeJSCode(sourceDir, code.CodeContent, logWriter)
	}

	if code.Runtime == "python" {
		return preparePythonCode(sourceDir, code.CodeContent, logWriter)
	}

	if code.Runtime == "go" {
		return prepareGoCode(sourceDir, code.CodeContent, logWriter)
	}

	var filename string
	switch code.Runtime {
	default:
		filename = "main.txt"
	}

	filePath := filepath.Join(sourceDir, filename)
	if err := os.WriteFile(filePath, []byte(code.CodeContent), 0644); err != nil {
		return fmt.Errorf("could not write serverless code to file: %w", err)
	}

	if logWriter != nil {
		fmt.Fprintf(logWriter, "📝 [Deployer] Wrote serverless function code to %s\n", filePath)
	}
	return nil
}

func prepareNodeJSCode(sourceDir, codeContent string, logWriter io.Writer) error {
	handlerPath := filepath.Join(sourceDir, "handler.js")
	if err := os.WriteFile(handlerPath, []byte(codeContent), 0644); err != nil {
		return err
	}

	pkgJSON := `{
  "name": "codedock-function",
  "version": "1.0.0",
  "main": "server.js",
  "scripts": {
    "start": "node server.js"
  },
  "dependencies": {
    "express": "^4.18.2"
  }
}`
	if err := os.WriteFile(filepath.Join(sourceDir, "package.json"), []byte(pkgJSON), 0644); err != nil {
		return err
	}

	serverJS := `const express = require('express');
const app = express();
app.use(express.json());

let handler;
try {
  handler = require('./handler.js');
} catch (err) {
  console.error("Failed to load handler.js:", err);
  process.exit(1);
}

app.all('*', async (req, res) => {
  try {
    if (typeof handler === 'function') {
      await handler(req, res);
    } else if (handler.default && typeof handler.default === 'function') {
      await handler.default(req, res);
    } else {
      res.status(500).send("No valid function exported in handler.js");
    }
  } catch (err) {
    console.error(err);
    res.status(500).send("Internal Server Error");
  }
});

const port = process.env.PORT || 3000;
app.listen(port, () => console.log('Function listening on port', port));`
	if err := os.WriteFile(filepath.Join(sourceDir, "server.js"), []byte(serverJS), 0644); err != nil {
		return err
	}

	if logWriter != nil {
		fmt.Fprintf(logWriter, "📝 [Deployer] Wrapped NodeJS serverless function with Express\n")
	}
	return nil
}

func preparePythonCode(sourceDir, codeContent string, logWriter io.Writer) error {
	handlerPath := filepath.Join(sourceDir, "handler.py")
	if err := os.WriteFile(handlerPath, []byte(codeContent), 0644); err != nil {
		return err
	}

	reqs := "Flask==2.3.2\nWerkzeug==2.3.4"
	if err := os.WriteFile(filepath.Join(sourceDir, "requirements.txt"), []byte(reqs), 0644); err != nil {
		return err
	}

	serverPy := `import os
from flask import Flask, request
import handler

app = Flask(__name__)

@app.route('/', defaults={'path': ''}, methods=['GET', 'POST', 'PUT', 'DELETE', 'PATCH'])
@app.route('/<path:path>', methods=['GET', 'POST', 'PUT', 'DELETE', 'PATCH'])
def catch_all(path):
    if hasattr(handler, 'main'):
        return handler.main(request)
    elif hasattr(handler, 'handler'):
        return handler.handler(request)
    else:
        return "No valid function exported in handler.py", 500

if __name__ == '__main__':
    port = int(os.environ.get('PORT', 3000))
    app.run(host='0.0.0.0', port=port)
`
	if err := os.WriteFile(filepath.Join(sourceDir, "main.py"), []byte(serverPy), 0644); err != nil {
		return err
	}

	if logWriter != nil {
		fmt.Fprintf(logWriter, "📝 [Deployer] Wrapped Python serverless function with Flask\n")
	}
	return nil
}

func prepareGoCode(sourceDir, codeContent string, logWriter io.Writer) error {
	handlerPath := filepath.Join(sourceDir, "handler.go")
	if err := os.WriteFile(handlerPath, []byte(codeContent), 0644); err != nil {
		return err
	}

	serverGo := `package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", Handler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	fmt.Printf("Listening on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}
`
	if err := os.WriteFile(filepath.Join(sourceDir, "main.go"), []byte(serverGo), 0644); err != nil {
		return err
	}

	goMod := `module codedock-function

go 1.23
`
	if err := os.WriteFile(filepath.Join(sourceDir, "go.mod"), []byte(goMod), 0644); err != nil {
		return err
	}

	if logWriter != nil {
		fmt.Fprintf(logWriter, "📝 [Deployer] Wrapped Go serverless function with net/http\n")
	}
	return nil
}
