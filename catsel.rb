class Catsel < Formula
  desc "Interactive file selector for concatenating and exporting text files"
  homepage "https://github.com/alexaldearroyo/catselector"
  url "https://github.com/alexaldearroyo/catselector/archive/refs/tags/v1.0.1.tar.gz"
  sha256 "579dfec00e4395b8d7160ef53bd87a109d4b926089cfba666965708ff4234c87"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args, "-o", bin/"catsel", "./"
  end

  test do
    system "#{bin}/catsel", "--version"
  end
end
