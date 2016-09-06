# system.d install

1. Install the `oregon-go` binary to `/usr/local/bin`
2. Create the `/etc/oregon-go.toml` config file. See the [sample file](../config.sample.toml) for details.
3. Copy the `oregon-go` systemd file to `/etc/default`
4. Copy the `oregon-go.service` systemd file to `/etc/systemd/system`
5. Then enable and run the service:

  ```sh
  systemctl daemon-reload
  systemctl enable oregon-go
  systemctl start oregon-go
  ```

You can check the status with `systemctl status oregon-go`.

Run `journalctl -u oregon-go` to see the logs.
