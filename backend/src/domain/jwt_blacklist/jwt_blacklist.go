package jwt_blacklist

type IJwtBlacklistService interface {
	AddToBlacklist(jwtToken string) error
	IsJwtInBlacklist(token string) (bool, error)
}
