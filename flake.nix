{
  description = "Simple flake";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-23.05";
  };

  outputs = {nixpkgs, ...}: let
    system = "x86_64-linux";
    pname = "hostsfile-discover";
    pkgs = import nixpkgs {inherit system;};
  in rec {
    packages.${system} = {
      default = pkgs.buildGoModule {
        name = pname;
        src = ./.;
        vendorSha256 = "sha256-pQpattmS9VmO3ZIQUFn66az8GSmB4IvYhTTCFn6SUmo=";
      };
    };
    nixosModules.default = {
      config,
      lib,
      ...
    }:
      with lib; let
        cfg = config.hostsfile-discover;
      in {
        options.hostsfile-discover = {
          enable = mkEnableOption "Enables the hostsfile-discover service";
          tld = mkOption {
            type = types.str;
            example = "home";
            default = "home";
            description = "Top level domain in LAN";
          };
          port = mkOption {
            type = types.int;
            example = 8080;
            default = 8080;
            description = "Port to listen on";
          };
        };

        config = mkIf cfg.enable {
          systemd.services.hostsfile-discover = {
            #wantedBy = ["multi-user.target"];
            after = ["network.target"];
            environment = {
              TLD = "${cfg.tld}";
              PORT = "${toString cfg.port}";
              HOSTS_FILE_PATH = "/etc/hosts";
            };
            serviceConfig = let
              selfpackage = packages.${system}.default;
            in {
              Restart = "on-failure";
              ExecStart = "${selfpackage}/bin/hostsfile-discover";
              DynamicUser = "yes";
              RuntimeDirectory = "hostsfile-discover";
              RuntimeDirectoryMode = "0755";
              StateDirectory = "hostsfile-discover";
              StateDirectoryMode = "0700";
              CacheDirectory = "hostsfile-discover";
              CacheDirectoryMode = "0750";
            };
          };
        };
      };
    formatter.${system} = pkgs.alejandra;
  };
}
