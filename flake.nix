{
  outputs = { nixpkgs, ... }: let
    system = "x86_64-linux";
    pkgs = nixpkgs.legacyPackages.${system};
  in rec {
    /*
    defaultPackage.${system} = python3WithPkgs.pkgs.buildPythonApplication {
      pname = "flake_ci";
      version = "0.0.0";
      src = ./.;
    };
    */

    devShell.${system} = (pkgs.mkShell {
      packages = with pkgs; [
        go
      ];
    });

    /*
    hydraJobs.test = defaultPackage;
    */
  };
}
