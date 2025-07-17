.PHONY: all test test-main test-tests test-safe clean lint build run fmt help benchmark coverage tidy dev-setup test-quick

# 기본 타겟
all: fmt lint test

# 도움말
help:
	@echo "사용 가능한 명령어:"
	@echo "  make all        - 포맷팅, 린팅, 테스트 실행"
	@echo "  make test       - 모든 테스트 실행"
	@echo "  make test-main  - 메인 프로젝트 테스트만 실행"
	@echo "  make test-tests - tests 디렉토리 테스트만 실행"
	@echo "  make test-safe  - 안전한 테스트만 실행 (동시성 제외)"
	@echo "  make lint       - 코드 린팅"
	@echo "  make fmt        - 코드 포맷팅"
	@echo "  make build      - 프로젝트 빌드"
	@echo "  make clean      - 테스트 캐시 정리"
	@echo "  make run        - 프로젝트 실행"

# 코드 포맷팅
fmt:
	@echo "코드 포맷팅..."
	@go fmt ./...
	@cd tests && go fmt ./...

# 모든 테스트 실행
test: test-main test-tests

# 메인 프로젝트 테스트
test-main:
	@echo "메인 프로젝트 테스트 실행 중..."
	@go test -v ./...

# tests 디렉토리 테스트
test-tests:
	@echo "tests 디렉토리 테스트 실행 중..."
	@cd tests && go test -v -timeout=30s ./... || (echo "일부 테스트가 실패했습니다."; exit 1)

# 안전한 테스트만 실행 (동시성 테스트 제외)
test-safe:
	@echo "안전한 테스트 실행 중..."
	@go test -v ./...
	@echo "동시성 테스트는 제외하고 실행합니다..."
	@cd tests && go test -v -run "TestTreeStructure|TestTreePath|TestTreeConsistency|TestNewTree|TestIsWildCard|TestIsCatchAll|TestSplitPath|TestMatch|TestInsertHandler|TestInsertChild|TestSetHandler|TestSplitNode|TestTryMatch" ./...

# 코드 린팅
lint:
	@echo "코드 린팅 중..."
	@golangci-lint run
	@cd tests && golangci-lint run

# 프로젝트 빌드
build:
	@echo "프로젝트 빌드 중..."
	@go build -o bin/liteframe ./...

# 프로젝트 실행
run:
	@echo "프로젝트 실행 중..."
	@go run ./...

# 테스트 캐시 정리
clean:
	@echo "테스트 캐시 정리 중..."
	@go clean -testcache
	@cd tests && go clean -testcache

# 벤치마크 테스트
benchmark:
	@echo "벤치마크 테스트 실행 중..."
	@cd tests && go test -bench=. -benchmem

# 테스트 커버리지
coverage:
	@echo "테스트 커버리지 생성 중..."
	@go test -coverprofile=coverage.out ./...
	@cd tests && go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "커버리지 리포트가 coverage.html에 생성되었습니다."

# 의존성 정리
tidy:
	@echo "의존성 정리 중..."
	@go mod tidy
	@cd tests && go mod tidy

# 개발 환경 설정
dev-setup:
	@echo "개발 환경 설정 중..."
	@go mod download
	@cd tests && go mod download
	@echo "개발 환경 설정이 완료되었습니다."

# 빠른 테스트 (상세 출력 없음)
test-quick:
	@echo "빠른 테스트 실행 중..."
	@go test ./...
	@cd tests && go test ./...