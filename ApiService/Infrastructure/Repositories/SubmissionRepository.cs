using Domain.Submissions;
using Infrastructure.Mappers;
using Microsoft.EntityFrameworkCore;
using Ports.OutBound.DTO;

namespace Infrastructure.Repositories;

public class SubmissionRepository
{
    private readonly AppDbContext _context;

    public SubmissionRepository(AppDbContext context)
    {
        _context = context;
    }

    public SubmissionResult GetSubmission(Guid submissionId)
    {
        try
        {
            var entity = _context.Submissions
                .AsNoTracking()
                .FirstOrDefault(x => x.Id == submissionId);

            return entity is null
                ? new SubmissionResult.Failure("submission_not_found")
                : new SubmissionResult.Success(SubmissionMapper.ToDomain(entity));
        }
        catch (Exception exception)
        {
            return new SubmissionResult.Failure(exception.Message);
        }
    }

    public SubmissionsResult GetSubmissionsByLabId(Guid labId)
    {
        try
        {
            var submissions = _context.Submissions
                .AsNoTracking()
                .Where(x => x.LabId == labId)
                .OrderByDescending(x => x.SubmittedAt)
                .Select(SubmissionMapper.ToDomain)
                .ToList();

            return new SubmissionsResult.Success(submissions);
        }
        catch (Exception exception)
        {
            return new SubmissionsResult.Failure(exception.Message);
        }
    }

    public SubmissionsResult GetSubmissions()
    {
        try
        {
            var submissions = _context.Submissions
                .AsNoTracking()
                .OrderByDescending(x => x.SubmittedAt)
                .Select(SubmissionMapper.ToDomain)
                .ToList();

            return new SubmissionsResult.Success(submissions);
        }
        catch (Exception exception)
        {
            return new SubmissionsResult.Failure(exception.Message);
        }
    }

    public SubmissionsResult GetSubmissionsByUserId(Guid userId)
    {
        try
        {
            var submissions = _context.Submissions
                .AsNoTracking()
                .Where(x => x.StudentId == userId)
                .OrderByDescending(x => x.SubmittedAt)
                .Select(SubmissionMapper.ToDomain)
                .ToList();

            return new SubmissionsResult.Success(submissions);
        }
        catch (Exception exception)
        {
            return new SubmissionsResult.Failure(exception.Message);
        }
    }

    public ParsedSubmissionsResult GetParsedSubmissions()
    {
        try
        {
            var parsedSubmissions = _context.ParsedSubmissions
                .AsNoTracking()
                .OrderByDescending(x => x.ParsedAt)
                .Select(ParsedSubmissionMapper.ToDomain)
                .ToList();

            return new ParsedSubmissionsResult.Success(parsedSubmissions);
        }
        catch (Exception exception)
        {
            return new ParsedSubmissionsResult.Failure(exception.Message);
        }
    }

    public SubmissionResult AddSubmission(Submission submission)
    {
        try
        {
            var entity = SubmissionMapper.ToEntity(submission);
            _context.Submissions.Add(entity);
            _context.SaveChanges();

            return new SubmissionResult.Success(SubmissionMapper.ToDomain(entity));
        }
        catch (Exception exception)
        {
            return new SubmissionResult.Failure(exception.Message);
        }
    }

    public ParsedSubmissionsResult GetParsedSubmissionByLabId(Guid labId)
    {
        try
        {
            var submissionIds = _context.Submissions
                .AsNoTracking()
                .Where(x => x.LabId == labId)
                .Select(x => x.Id);

            var parsedSubmissions = _context.ParsedSubmissions
                .AsNoTracking()
                .Where(x => submissionIds.Contains(x.SubmissionId))
                .OrderByDescending(x => x.ParsedAt)
                .Select(ParsedSubmissionMapper.ToDomain)
                .ToList();

            return new ParsedSubmissionsResult.Success(parsedSubmissions);
        }
        catch (Exception exception)
        {
            return new ParsedSubmissionsResult.Failure(exception.Message);
        }
    }
}
