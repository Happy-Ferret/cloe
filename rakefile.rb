TOTAL_COVERAGE_FILE = 'coverage.txt'.freeze # This path is specified by codecov.
BIN_PATH = File.absolute_path 'bin'

task :deps do
  sh %w[go get -u
        github.com/alecthomas/gometalinter
        github.com/mattn/goveralls
        github.com/raviqqe/liche].join ' '
  sh 'gometalinter --install'
  sh 'go get -d -t ./...'
  sh 'gem install rake rubocop'
end

task :build do
  sh 'go build -o bin/tisp src/cmd/tisp/main.go'
end

task :fast_unit_test do
  sh 'go test ./...'
end

task :unit_test do
  coverage_file = "/tmp/tisp-unit-test-#{Process.pid}.coverage"

  sh "echo mode: atomic > #{TOTAL_COVERAGE_FILE}"

  `go list ./src/lib/...`.split.each do |package|
    sh %W[go test
          -covermode atomic
          -coverprofile #{coverage_file}
          #{`uname -m` =~ /x86_64/ ? '-race' : ''}
          #{package}].join ' '

    verbose false do
      if File.exist? coverage_file
        sh "cat #{coverage_file} | grep -v mode: >> #{TOTAL_COVERAGE_FILE}"
        rm coverage_file
      end
    end
  end
end

task command_test: :build do
  sh 'bundler install'
  sh %W[bundler exec cucumber
        -r examples/aruba.rb
        PATH=#{BIN_PATH}:$PATH
        examples].join ' '
end

task test: %i[unit_test command_test]

task :format do
  sh 'go fix ./...'
  sh 'go fmt ./...'

  Dir.glob '**/*.go' do |file|
    sh "goimports -w #{file}"
  end

  sh 'rubocop -a'
end

task :lint do
  sh %w[gometalinter
        --disable gocyclo
        --disable vetshadow
        --enable gofmt
        --enable goimports
        --enable misspell
        ./...].join ' '
  sh 'rubocop'
  sh "liche -v #{Dir.glob('**/*.md').join ' '}"
end

task install: %i[deps test build] do
  sh 'go get ./...'
end

task default: %i[test build]

task :clean do
  sh 'git clean -dfx'
end
