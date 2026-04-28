using Domain.Users;

namespace Presentation.ResponseDTO.Auths;

public sealed record RegisterResponse(
        
    Guid UserId,
    
    string Email,
    
    string FullName,
    
    string? GroupName,
    
    UserRole Role);