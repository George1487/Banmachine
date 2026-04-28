using Domain.Users;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Ports.InBound.Contracts;
using Ports.InBound.DTO;
using Presentation.RequestDTO;
using Presentation.ResponseDTO.Auths;
using RegisterRequest = Presentation.RequestDTO.RegisterRequest;

namespace Presentation.Controllers;

[ApiController]
[Route("api/auth")]
public class AuthenticationController : ControllerBase
{
    private readonly ITokenService _tokenService;
    
    private readonly IUserService _userService;
    
    public AuthenticationController(ITokenService tokenService,
        IUserService userService)
    {
        _tokenService =  tokenService;
        _userService = userService;
    }
    
    
    [HttpPost("register")]
    public ActionResult<RegisterResponse> Register([FromBody] RegisterRequest request)
    {
        var userResult = _userService.RegisterUser(
            new User(Guid.NewGuid(),
                request.Email,
                request.FullName,
                request.Password,
                request.GroupName,
                request.Role));

        if (userResult is UserResult.Failure failure)
        {
            return BadRequest(failure.Reason);
        }
        
        var success = (UserResult.Success)userResult;
        return Ok(new RegisterResponse(
            success.User.UserId,
            success.User.Email,
            success.User.FullName,
            success.User.GroupName,
            success.User.Role));
    }

    [HttpPost("login")]
    public ActionResult<LoginResponse> Login([FromBody] LoginRequest request)
    {
        LoginResult result = _userService.Login(request.Email, request.Password);
        if (result is LoginResult.Failure failure)
        {
            return BadRequest(failure.Reason);
        }
        
        LoginResult.Success success = (LoginResult.Success)result;
        
        return Ok(new LoginResponse(success.Token, 
            new LoginUserResponse(success.User.UserId, 
                success.User.FullName, success.User.Role)));
    }
    
    [Authorize]
    [HttpGet("me")]
    public ActionResult<GetMeResponse> GetMe()
    {
        var stringId = HttpContext.User.Claims
            .FirstOrDefault(c => c.Type == "UserId")!.Value;
        
        var userId = Guid.Parse(stringId);
        
        UserResult userResult = _userService.GetUser(userId);
        if (userResult is UserResult.Failure failure)
        {
            return BadRequest(failure.Reason);
        }
        
        UserResult.Success success = (UserResult.Success)userResult;
        return Ok(new GetMeResponse(
            success.User.UserId,
            success.User.Email, 
            success.User.FullName,
            success.User.Role));
    }
}