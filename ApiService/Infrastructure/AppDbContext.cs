using Infrastructure.Entities;
using Microsoft.EntityFrameworkCore;

namespace Infrastructure;

public class AppDbContext : DbContext
{
    public DbSet<UserEntity> Users => Set<UserEntity>();
    public DbSet<LabEntity> Labs => Set<LabEntity>();
    public DbSet<SubmissionEntity> Submissions => Set<SubmissionEntity>();
    public DbSet<ParsedSubmissionEntity> ParsedSubmissions => Set<ParsedSubmissionEntity>();
    public DbSet<IngestJobEntity> IngestJobs => Set<IngestJobEntity>();
    public DbSet<AnalysisJobEntity> AnalysisJobs => Set<AnalysisJobEntity>();
    public DbSet<PairwiseSimilarityEntity> PairwiseSimilarities => Set<PairwiseSimilarityEntity>();
    public DbSet<SubmissionAnalysisSummaryEntity> SubmissionAnalysisSummaries => Set<SubmissionAnalysisSummaryEntity>();

    public AppDbContext(DbContextOptions<AppDbContext> options) : base(options)
    {
    }

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        base.OnModelCreating(modelBuilder);

        modelBuilder.Entity<UserEntity>(entity =>
        {
            entity.ToTable("users");
            entity.HasKey(x => x.Id);

            entity.Property(x => x.Id).HasColumnName("user_id");
            entity.Property(x => x.Email).HasColumnName("email");
            entity.Property(x => x.Password).HasColumnName("password");
            entity.Property(x => x.FullName).HasColumnName("full_name");
            entity.Property(x => x.Role).HasColumnName("role");
            entity.Property(x => x.GroupName).HasColumnName("group_name");

            entity.HasIndex(x => x.Email).IsUnique();
        });

        modelBuilder.Entity<LabEntity>(entity =>
        {
            entity.ToTable("labs");
            entity.HasKey(x => x.Id);

            entity.Property(x => x.Id).HasColumnName("lab_id");
            entity.Property(x => x.TeacherId).HasColumnName("teacher_id");
            entity.Property(x => x.Title).HasColumnName("title");
            entity.Property(x => x.Description).HasColumnName("description");
            entity.Property(x => x.Status).HasColumnName("status");
            entity.Property(x => x.DeadlineAt).HasColumnName("deadline_at");
        });

        modelBuilder.Entity<SubmissionEntity>(entity =>
        {
            entity.ToTable("submissions");
            entity.HasKey(x => x.Id);

            entity.Property(x => x.Id).HasColumnName("submission_id");
            entity.Property(x => x.LabId).HasColumnName("lab_id");
            entity.Property(x => x.StudentId).HasColumnName("student_id");
            entity.Property(x => x.Status).HasColumnName("status");
            entity.Property(x => x.SourceFileName).HasColumnName("source_file_name");
            entity.Property(x => x.MimeType).HasColumnName("mime_type");
            entity.Property(x => x.StorageKey).HasColumnName("storage_key");
            entity.Property(x => x.SubmittedAt).HasColumnName("submitted_at");
        });

        modelBuilder.Entity<ParsedSubmissionEntity>(entity =>
        {
            entity.ToTable("parsed_submissions");
            entity.HasKey(x => x.Id);

            entity.Property(x => x.Id).HasColumnName("parsed_submission_id");
            entity.Property(x => x.SubmissionId).HasColumnName("submission_id");
            entity.Property(x => x.RawText).HasColumnName("raw_text");
            entity.Property(x => x.StructuredData)
                .HasColumnName("structured_data")
                .HasColumnType("jsonb");
            entity.Property(x => x.ParsedAt).HasColumnName("parsed_at");
        });

        modelBuilder.Entity<IngestJobEntity>(entity =>
        {
            entity.ToTable("ingest_jobs");
            entity.HasKey(x => x.Id);

            entity.Property(x => x.Id).HasColumnName("ingest_job_id");
            entity.Property(x => x.SubmissionId).HasColumnName("submission_id");
            entity.Property(x => x.Status).HasColumnName("status");
            entity.Property(x => x.CreatedAt).HasColumnName("created_at");
            entity.Property(x => x.FinishedAt).HasColumnName("finished_at");
            entity.Property(x => x.ErrorMessage).HasColumnName("error_message");
        });

        modelBuilder.Entity<AnalysisJobEntity>(entity =>
        {
            entity.ToTable("analysis_jobs");
            entity.HasKey(x => x.Id);

            entity.Property(x => x.Id).HasColumnName("analysis_job_id");
            entity.Property(x => x.LabId).HasColumnName("lab_id");
            entity.Property(x => x.Status).HasColumnName("status");
            entity.Property(x => x.CreatedBy).HasColumnName("created_by");
            entity.Property(x => x.CreatedAt).HasColumnName("created_at");
            entity.Property(x => x.StartedAt).HasColumnName("started_at");
            entity.Property(x => x.FinishedAt).HasColumnName("finished_at");
            entity.Property(x => x.ErrorMessage).HasColumnName("error_message");
        });

        modelBuilder.Entity<PairwiseSimilarityEntity>(entity =>
        {
            entity.ToTable("pairwise_similarities");
            entity.HasKey(x => x.Id);

            entity.Property(x => x.Id).HasColumnName("pairwise_similarity_id");
            entity.Property(x => x.AnalysisJobId).HasColumnName("analysis_job_id");
            entity.Property(x => x.LabId).HasColumnName("lab_id");
            entity.Property(x => x.LeftSubmissionId).HasColumnName("left_submission_id");
            entity.Property(x => x.RightSubmissionId).HasColumnName("right_submission_id");
            entity.Property(x => x.TextScore).HasColumnName("text_score");
            entity.Property(x => x.CalculationScore).HasColumnName("calculation_score");
            entity.Property(x => x.ImagesScore).HasColumnName("images_score");
            entity.Property(x => x.FinalScore).HasColumnName("final_score");
        });

        modelBuilder.Entity<SubmissionAnalysisSummaryEntity>(entity =>
        {
            entity.ToTable("submission_analysis_summaries");
            entity.HasKey(x => x.Id);

            entity.Property(x => x.Id).HasColumnName("submission_analysis_summary_id");
            entity.Property(x => x.AnalysisJobId).HasColumnName("analysis_job_id");
            entity.Property(x => x.SubmissionId).HasColumnName("submission_id");
            entity.Property(x => x.TopMatchSubmissionId).HasColumnName("top_match_submission_id");
            entity.Property(x => x.TopMatchScore).HasColumnName("top_match_score");
            entity.Property(x => x.FinalScoreRiskLevel).HasColumnName("final_score_risk_level");
        });
    }
}
