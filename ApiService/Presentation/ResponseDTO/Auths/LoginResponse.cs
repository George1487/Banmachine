namespace Presentation.ResponseDTO.Auths;

public sealed record LoginResponse(string Token, LoginUserResponse User);