{
  description = "Simple flake";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-22.11";
  };

  outputs = {nixpkgs, ...}: let
    system = "x86_64-linux";
    pname = "hostsfile-discover";
  in {
    packages.${system} = let
      pkgs = import nixpkgs {inherit system;}; 
    in {
      default = pkgs.buildGoModule {
        name = pname;
        src = ./.;
        vendorSha256 = "sha256-pQpattmS9VmO3ZIQUFn66az8GSmB4IvYhTTCFn6SUmo=";
      };
    };
  };
}
