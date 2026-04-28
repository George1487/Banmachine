using Domain.Users;

namespace Presentation.ResponseDTO.Auths;

public sealed record LoginUserResponse(
    
    Guid UserId,
    
    string FullName,
    
    UserRole Role 
    );
    