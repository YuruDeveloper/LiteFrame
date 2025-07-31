.PHONY: all help clean fmt lint test test-quick test-safe benchmark pprof-gethandler coverage tidy dev-setup

# =============================================================================
# 기본 타겟
# =============================================================================
all: fmt lint test

help:
	@echo "=== LiteFrame Makefile 사용 가능한 명령어 ==="
	@echo ""
	@echo "개발 워크플로우:"
	@echo "  make all         - 포맷팅, 린팅, 테스트 실행 (권장)"
	@echo "  make dev-setup   - 개발 환경 초기 설정"
	@echo ""
	@echo "코드 품질:"
	@echo "  make fmt         - 코드 포맷팅 (Go fmt)"
	@echo "  make lint        - 코드 린팅 (golangci-lint)"
	@echo "  make tidy        - Go 모듈 의존성 정리"
	@echo ""
	@echo "테스트:"
	@echo "  make test        - 메인 프로젝트 테스트 실행"
	@echo "  make test-quick  - 빠른 테스트 (상세 출력 없음)"
	@echo "  make test-safe   - 안전한 테스트만 실행 (동시성 제외)"
	@echo ""
	@echo "성능 분석:"
	@echo "  make benchmark   - 벤치마크 테스트 실행"
	@echo "  make pprof-gethandler - GetHandler 성능 프로파일링"
	@echo "  make coverage    - 테스트 커버리지 리포트 생성"
	@echo ""
	@echo "유틸리티:"
	@echo "  make clean       - 테스트 캐시 정리"
	@echo ""

# =============================================================================
# 코드 품질 관리
# =============================================================================
fmt:
	@echo "📝 코드 포맷팅 중..."
	@go fmt ./...

lint:
	@echo "🔍 코드 린팅 중..."
	@golangci-lint run

tidy:
	@echo "📦 Go 모듈 의존성 정리 중..."
	@go mod tidy
	@cd bench && go mod tidy

# =============================================================================
# 테스트
# =============================================================================
test:
	@echo "🧪 메인 프로젝트 테스트 실행 중..."
	@go test -v ./...

test-quick:
	@echo "⚡ 빠른 테스트 실행 중..."
	@go test ./...

test-safe:
	@echo "🛡️ 안전한 테스트 실행 중 (동시성 테스트 제외)..."
	@go test -v -run "TestTreeStructure|TestTreeConsistency|TestNewTree|TestIsWildCard|TestIsCatchAll|TestPathWithSegment|TestMatch|TestInsertChild|TestSetHandler|TestSplitNode|TestGetHandler|TestComplexRouting|TestAdditionalScenarios|TestEdgeCases|TestInputValidation" ./...

# =============================================================================
# 성능 분석
# =============================================================================
benchmark:
	@echo "🏃 벤치마크 테스트 실행 중..."
	@cd bench && go test -bench=. -benchmem

pprof-gethandler:
	@echo "🔬 GetHandler 성능 프로파일링 중..."
	@cd bench && go test -bench=BenchmarkGetHandler -run=^$$ -benchmem -cpuprofile=CPU.prof && go tool pprof bench.test CPU.prof

coverage:
	@echo "📊 테스트 커버리지 생성 중..."
	@go test -coverprofile=main_coverage.out ./...
	@go tool cover -html=main_coverage.out -o main_coverage.html
	@echo "✅ 커버리지 리포트 생성 완료: main_coverage.html"

# =============================================================================
# 유틸리티
# =============================================================================
clean:
	@echo "🧹 테스트 캐시 정리 중..."
	@go clean -testcache
	@cd bench && go clean -testcache
	@rm -f *.out *.html CPU.prof bench/CPU.prof bench/*.test

dev-setup:
	@echo "🚀 개발 환경 설정 중..."
	@go mod download
	@cd bench && go mod download
	@echo "✅ 개발 환경 설정 완료"