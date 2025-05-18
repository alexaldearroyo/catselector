class Catsel < Formula
  desc "Interactive file selector for concatenating and exporting text files"
  homepage "https://github.com/alexaldearroyo/catselector"
  url "https://github.com/alexaldearroyo/catselector/archive/refs/tags/v1.0.3.tar.gz"
  sha256 "00daa38e0c4ab6865360ba74457baf2b85b42d7950598da7d2c0efddc95c23d5"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args, "-o", bin/"catsel", "./"
  end

  test do
    system "#{bin}/catsel", "--version"
  end
end
