ALTER TABLE server_settings ADD COLUMN site_name TEXT DEFAULT '';
ALTER TABLE server_settings ADD COLUMN public_ipv4 TEXT DEFAULT '';
ALTER TABLE server_settings ADD COLUMN public_ipv6 TEXT DEFAULT '';
ALTER TABLE server_settings ADD COLUMN show_sponsorship_popup BOOLEAN DEFAULT 1;
ALTER TABLE server_settings ADD COLUMN disable_two_step_confirmation BOOLEAN DEFAULT 0;
