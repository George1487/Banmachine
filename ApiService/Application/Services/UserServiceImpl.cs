using Domain.Users;
using Ports.InBound.Contracts;
using Ports.InBound.DTO;
using Ports.OutBound.Contracts;

namespace Application.Services;

public class UserServiceImpl : IUserService
{

    private readonly IUserRepository _userRepository;
    
    private readonly ITokenService _tokenService;
    
    public UserServiceImpl(IUserRepository userRepository,
        ITokenService tokenService)
    {
        _userRepository = userRepository;
        _tokenService = tokenService;
    }
    
    public UserResult RegisterUser(User user)
    {
        return _userRepository.AddUser(user);
    }

    public LoginResult Login(string email, string password)
    {
        // if (email == "zvo" && password == "zvo")
        // {
        //     return _tokenService.CreateToken(new User(
        //         Guid.NewGuid(),email, "xuy", password, null, UserType.Teacher ));
        // }
        
        var user = _userRepository.GetUserByEmail(email);

        if (user is UserResult.Failure failure)
        {
            return new LoginResult.Failure(failure.Reason);
        }
        
        var success = (UserResult.Success)user;
        
        if (success.User.Password != password)
        {
            return new LoginResult.Failure("Wrong password");
        }

        return new LoginResult.Success(_tokenService.CreateToken(success.User), success.User);
    }

    public UserResult GetUser(Guid userId)
    {
        return _userRepository.GetUser(userId);
    }
}