using Domain.Users;

namespace Presentation.ResponseDTO.Auths;

public sealed record GetMeResponse(
        
    Guid UserId,
    
    string Email,
    
    string FullName,
    
    UserRole Role );