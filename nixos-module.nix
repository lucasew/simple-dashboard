{config, pkgs, lib, ...}:
let
  cfg = config.services.simple-dashboardd;
in {
  options = with lib; {
    services.simple-dashboardd = {
      enable = mkEnableOption "Webapp to show system usage";
      package = mkOption {
        description = "Simple dashboard package";
        type = lib.types.package;
        default = pkgs.callPackage ./package.nix {};
      };
      config = mkOption {
        description = "Config string to be used by simple-dashboard";
        type = types.lines;
        default = builtins.readFile ./config.ini.example;
      };
      port = mkOption {
        description = "Port to listen";
        default = 42069;
        type = types.port;
      };
      openFirewall = mkEnableOption "Open the dashboard port on the firewall (recommended)";
    };
  };
  config = with lib; mkIf cfg.enable {
    systemd.services.simple-dashboardd = {
      description = "Simple system usage dashboard";
      # Without wantedBy, enable=true only defines a unit that never joins the boot graph.
      wantedBy = [ "multi-user.target" ];
      after = [ "network.target" ];
      script = ''
        ${lib.getExe cfg.package} -c ${builtins.toFile "simple-dashboard.cfg" cfg.config} -p ${builtins.toString cfg.port}
      '';
      serviceConfig = {
        Restart = "on-failure";
        RestartSec = "5s";
      };
      restartIfChanged = true;
    };
    networking.firewall.allowedTCPPorts = mkIf cfg.openFirewall [ cfg.port ];
  };
}
