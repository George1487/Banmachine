using Domain.Users;

namespace Ports.InBound.DTO;
public abstract record LoginResult
{
    
    private LoginResult() {}
    
    public sealed record Failure(string Reason) : LoginResult;
    
    public sealed record Success(string Token, User User) : LoginResult;
    
}