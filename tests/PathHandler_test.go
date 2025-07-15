package tests

import (
	"testing"
	"LiteFrame/Router/Tree/Component"
	"github.com/stretchr/testify/assert"
)

func TestNewPathHandler(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	pathHandler := Component.NewPathHandler(err, "/api/users")
	
	assert.NotNil(t, pathHandler)
	assert.Equal(t, "/api/users", pathHandler.Path)
}

func TestPathHandler_GetPath(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	pathHandler := Component.NewPathHandler(err, "/api/users")
	
	assert.Equal(t, "/api/users", pathHandler.GetPath())
}

func TestPathHandler_SetPath(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	pathHandler := Component.NewPathHandler(err, "/api/users")
	
	// 성공적인 경로 설정
	err = pathHandler.SetPath("/api/posts")
	assert.NoError(t, err)
	assert.Equal(t, "/api/posts", pathHandler.GetPath())
	
	// 빈 경로 설정 시도
	err = pathHandler.SetPath("")
	assert.Error(t, err)
	assert.Equal(t, "/api/posts", pathHandler.GetPath()) // 기존 경로 유지
}

func TestPathHandler_Match_ExactMatch(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	pathHandler := Component.NewPathHandler(err, "/api/users")
	
	// 정확히 일치하는 경우
	matched, matchingChar, leftPath := pathHandler.Match("/api/users")
	assert.True(t, matched)
	assert.Equal(t, 10, matchingChar) // "/api/users" 길이
	assert.Equal(t, "", leftPath)
}

func TestPathHandler_Match_PartialMatch(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	pathHandler := Component.NewPathHandler(err, "/api/users")
	
	// 부분적으로 일치하는 경우
	matched, matchingChar, leftPath := pathHandler.Match("/api/users/123")
	assert.True(t, matched)
	assert.Equal(t, 10, matchingChar) // "/api/users" 길이
	assert.Equal(t, "/123", leftPath)
}

func TestPathHandler_Match_NoMatch(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	pathHandler := Component.NewPathHandler(err, "/api/users")
	
	// 전혀 일치하지 않는 경우
	matched, matchingChar, leftPath := pathHandler.Match("/api/posts")
	assert.False(t, matched)
	assert.Equal(t, 5, matchingChar) // "/api/" 까지만 일치
	assert.Equal(t, "posts", leftPath)
}

func TestPathHandler_Match_ShorterPath(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	pathHandler := Component.NewPathHandler(err, "/api/users")
	
	// 더 짧은 경로와 매칭
	matched, matchingChar, leftPath := pathHandler.Match("/api")
	assert.False(t, matched)
	assert.Equal(t, 4, matchingChar) // "/api" 길이
	assert.Equal(t, "", leftPath)
}

func TestPathHandler_Match_LongerHandlerPath(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	pathHandler := Component.NewPathHandler(err, "/api/users/profile")
	
	// 핸들러 경로가 더 긴 경우
	matched, matchingChar, leftPath := pathHandler.Match("/api/users")
	assert.False(t, matched)
	assert.Equal(t, 10, matchingChar) // "/api/users" 길이
	assert.Equal(t, "", leftPath)
}

func TestPathHandler_Match_EmptyPaths(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	pathHandler := Component.NewPathHandler(err, "")
	
	// 빈 경로 테스트
	matched, matchingChar, leftPath := pathHandler.Match("")
	assert.True(t, matched)
	assert.Equal(t, 0, matchingChar)
	assert.Equal(t, "", leftPath)
	
	// 핸들러 경로가 빈 경우
	matched, matchingChar, leftPath = pathHandler.Match("/api")
	assert.True(t, matched)
	assert.Equal(t, 0, matchingChar)
	assert.Equal(t, "/api", leftPath)
}

func TestPathHandler_Match_SpecialCharacters(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	pathHandler := Component.NewPathHandler(err, "/api/users-123")
	
	// 특수 문자 포함 경로
	matched, matchingChar, leftPath := pathHandler.Match("/api/users-123/profile")
	assert.True(t, matched)
	assert.Equal(t, 14, matchingChar) // "/api/users-123" 길이
	assert.Equal(t, "/profile", leftPath)
}

func TestPathHandler_Match_UnicodeCharacters(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	pathHandler := Component.NewPathHandler(err, "/api/사용자")
	
	// 유니코드 문자 포함 경로
	matched, matchingChar, leftPath := pathHandler.Match("/api/사용자/프로필")
	assert.True(t, matched)
	assert.Equal(t, 14, matchingChar) // "/api/사용자" 길이 (UTF-8 바이트 기준)
	assert.Equal(t, "/프로필", leftPath)
}

func TestPathHandler_Match_EdgeCases(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	pathHandler := Component.NewPathHandler(err, "/")
	
	// 루트 경로 테스트
	matched, matchingChar, leftPath := pathHandler.Match("/api")
	assert.True(t, matched)
	assert.Equal(t, 1, matchingChar) // "/" 길이
	assert.Equal(t, "api", leftPath)
	
	// 완전히 일치하는 루트 경로
	matched, matchingChar, leftPath = pathHandler.Match("/")
	assert.True(t, matched)
	assert.Equal(t, 1, matchingChar)
	assert.Equal(t, "", leftPath)
}

func TestPathHandler_Match_ConsecutiveSlashes(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	pathHandler := Component.NewPathHandler(err, "/api//users")
	
	// 연속된 슬래시 테스트
	matched, matchingChar, leftPath := pathHandler.Match("/api//users/123")
	assert.True(t, matched)
	assert.Equal(t, 11, matchingChar) // "/api//users" 길이 = 11
	assert.Equal(t, "/123", leftPath)
}

func TestPathHandler_Match_CaseSensitive(t *testing.T) {
	err := Component.NewError(Component.StaticType, "테스트 에러", "/test")
	pathHandler := Component.NewPathHandler(err, "/API/Users")
	
	// 대소문자 구분 테스트
	matched, matchingChar, leftPath := pathHandler.Match("/api/users")
	assert.False(t, matched)
	assert.Equal(t, 1, matchingChar) // "/" 까지만 일치
	assert.Equal(t, "api/users", leftPath)
}