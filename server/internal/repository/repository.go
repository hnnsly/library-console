package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hnnsly/library-console/internal/repository/postgres"
	"github.com/hnnsly/library-console/internal/repository/redis"
)

// LibraryRepository фасадный репозиторий для работы с БД и кешем
type LibraryRepository struct {
	pg *postgres.Queries
	rd *redis.Redis
}

// New создает новый LibraryRepository
func New(pg *postgres.Queries, rd *redis.Redis) *LibraryRepository {
	return &LibraryRepository{
		pg: pg,
		rd: rd,
	}
}

// ===== Analytics Methods =====

func (r *LibraryRepository) GetActiveReaderStatistics(ctx context.Context) ([]*postgres.GetActiveReaderStatisticsRow, error) {
	return r.pg.GetActiveReaderStatistics(ctx)
}

func (r *LibraryRepository) GetBooksByAuthorInHall(ctx context.Context, arg postgres.GetBooksByAuthorInHallParams) ([]*postgres.GetBooksByAuthorInHallRow, error) {
	return r.pg.GetBooksByAuthorInHall(ctx, arg)
}

func (r *LibraryRepository) GetHallUtilizationReport(ctx context.Context) ([]*postgres.GetHallUtilizationReportRow, error) {
	return r.pg.GetHallUtilizationReport(ctx)
}

func (r *LibraryRepository) GetLibraryStatistics(ctx context.Context) (*postgres.GetLibraryStatisticsRow, error) {
	return r.pg.GetLibraryStatistics(ctx)
}

func (r *LibraryRepository) GetMonthlyStatistics(ctx context.Context, arg postgres.GetMonthlyStatisticsParams) ([]*postgres.GetMonthlyStatisticsRow, error) {
	return r.pg.GetMonthlyStatistics(ctx, arg)
}

func (r *LibraryRepository) GetPopularBooks(ctx context.Context, limitVal int32) ([]*postgres.GetPopularBooksRow, error) {
	return r.pg.GetPopularBooks(ctx, limitVal)
}

func (r *LibraryRepository) GetReadersWithSingleCopyBooks(ctx context.Context) ([]*postgres.GetReadersWithSingleCopyBooksRow, error) {
	return r.pg.GetReadersWithSingleCopyBooks(ctx)
}

// ===== Authors Methods =====

func (r *LibraryRepository) CountAuthors(ctx context.Context) (int64, error) {
	return r.pg.CountAuthors(ctx)
}

func (r *LibraryRepository) CreateAuthor(ctx context.Context, arg postgres.CreateAuthorParams) (*postgres.Author, error) {
	return r.pg.CreateAuthor(ctx, arg)
}

func (r *LibraryRepository) DeleteAuthor(ctx context.Context, authorID uuid.UUID) error {
	return r.pg.DeleteAuthor(ctx, authorID)
}

func (r *LibraryRepository) GetAuthorByID(ctx context.Context, authorID uuid.UUID) (*postgres.Author, error) {
	return r.pg.GetAuthorByID(ctx, authorID)
}

func (r *LibraryRepository) GetAuthorsByBook(ctx context.Context, bookID uuid.UUID) ([]*postgres.Author, error) {
	return r.pg.GetAuthorsByBook(ctx, bookID)
}

func (r *LibraryRepository) ListAuthors(ctx context.Context, arg postgres.ListAuthorsParams) ([]*postgres.Author, error) {
	return r.pg.ListAuthors(ctx, arg)
}

func (r *LibraryRepository) SearchAuthorsByName(ctx context.Context, searchQuery string) ([]*postgres.Author, error) {
	return r.pg.SearchAuthorsByName(ctx, searchQuery)
}

func (r *LibraryRepository) UpdateAuthor(ctx context.Context, arg postgres.UpdateAuthorParams) (*postgres.Author, error) {
	return r.pg.UpdateAuthor(ctx, arg)
}

// ===== Book Authors Methods =====

func (r *LibraryRepository) AddBookAuthor(ctx context.Context, arg postgres.AddBookAuthorParams) error {
	return r.pg.AddBookAuthor(ctx, arg)
}

func (r *LibraryRepository) GetAuthorBooks(ctx context.Context, authorID uuid.UUID) ([]*postgres.Book, error) {
	return r.pg.GetAuthorBooks(ctx, authorID)
}

func (r *LibraryRepository) GetBookAuthors(ctx context.Context, bookID uuid.UUID) ([]*postgres.Author, error) {
	return r.pg.GetBookAuthors(ctx, bookID)
}

func (r *LibraryRepository) RemoveAllBookAuthors(ctx context.Context, bookID uuid.UUID) error {
	return r.pg.RemoveAllBookAuthors(ctx, bookID)
}

func (r *LibraryRepository) RemoveBookAuthor(ctx context.Context, arg postgres.RemoveBookAuthorParams) error {
	return r.pg.RemoveBookAuthor(ctx, arg)
}

// ===== Book Copies Methods =====

func (r *LibraryRepository) CountAvailableBookCopies(ctx context.Context, bookID uuid.UUID) (int64, error) {
	return r.pg.CountAvailableBookCopies(ctx, bookID)
}

func (r *LibraryRepository) CountBookCopiesByBook(ctx context.Context, bookID uuid.UUID) (int64, error) {
	return r.pg.CountBookCopiesByBook(ctx, bookID)
}

func (r *LibraryRepository) CreateBookCopy(ctx context.Context, arg postgres.CreateBookCopyParams) (*postgres.BookCopy, error) {
	return r.pg.CreateBookCopy(ctx, arg)
}

func (r *LibraryRepository) DeleteBookCopy(ctx context.Context, copyID uuid.UUID) error {
	return r.pg.DeleteBookCopy(ctx, copyID)
}

func (r *LibraryRepository) GetBookCopiesByStatus(ctx context.Context, status postgres.NullBookStatus) ([]*postgres.GetBookCopiesByStatusRow, error) {
	return r.pg.GetBookCopiesByStatus(ctx, status)
}

func (r *LibraryRepository) GetBookCopyByCode(ctx context.Context, copyCode string) (*postgres.GetBookCopyByCodeRow, error) {
	return r.pg.GetBookCopyByCode(ctx, copyCode)
}

func (r *LibraryRepository) GetBookCopyByID(ctx context.Context, copyID uuid.UUID) (*postgres.GetBookCopyByIDRow, error) {
	return r.pg.GetBookCopyByID(ctx, copyID)
}

func (r *LibraryRepository) ListAvailableBookCopies(ctx context.Context) ([]*postgres.ListAvailableBookCopiesRow, error) {
	return r.pg.ListAvailableBookCopies(ctx)
}

func (r *LibraryRepository) ListBookCopiesByBook(ctx context.Context, bookID uuid.UUID) ([]*postgres.ListBookCopiesByBookRow, error) {
	return r.pg.ListBookCopiesByBook(ctx, bookID)
}

func (r *LibraryRepository) UpdateBookCopy(ctx context.Context, arg postgres.UpdateBookCopyParams) (*postgres.BookCopy, error) {
	return r.pg.UpdateBookCopy(ctx, arg)
}

// ===== Book Issues Methods =====

func (r *LibraryRepository) CountActiveIssuesByReader(ctx context.Context, readerID uuid.UUID) (int64, error) {
	return r.pg.CountActiveIssuesByReader(ctx, readerID)
}

func (r *LibraryRepository) CountOverdueIssues(ctx context.Context) (int64, error) {
	return r.pg.CountOverdueIssues(ctx)
}

func (r *LibraryRepository) CreateBookIssue(ctx context.Context, arg postgres.CreateBookIssueParams) (*postgres.BookIssue, error) {
	return r.pg.CreateBookIssue(ctx, arg)
}

func (r *LibraryRepository) ExtendBookIssue(ctx context.Context, arg postgres.ExtendBookIssueParams) error {
	return r.pg.ExtendBookIssue(ctx, arg)
}

func (r *LibraryRepository) GetActiveIssuesByReader(ctx context.Context, readerID uuid.UUID) ([]*postgres.GetActiveIssuesByReaderRow, error) {
	return r.pg.GetActiveIssuesByReader(ctx, readerID)
}

func (r *LibraryRepository) GetAllActiveIssues(ctx context.Context) ([]*postgres.GetAllActiveIssuesRow, error) {
	return r.pg.GetAllActiveIssues(ctx)
}

func (r *LibraryRepository) GetBookIssueByID(ctx context.Context, issueID uuid.UUID) (*postgres.GetBookIssueByIDRow, error) {
	return r.pg.GetBookIssueByID(ctx, issueID)
}

func (r *LibraryRepository) GetBookIssueHistory(ctx context.Context, copyID uuid.UUID) ([]*postgres.GetBookIssueHistoryRow, error) {
	return r.pg.GetBookIssueHistory(ctx, copyID)
}

func (r *LibraryRepository) GetIssueDueSoon(ctx context.Context) ([]*postgres.GetIssueDueSoonRow, error) {
	return r.pg.GetIssueDueSoon(ctx)
}

func (r *LibraryRepository) GetIssueHistory(ctx context.Context, arg postgres.GetIssueHistoryParams) ([]*postgres.GetIssueHistoryRow, error) {
	return r.pg.GetIssueHistory(ctx, arg)
}

func (r *LibraryRepository) GetOverdueIssues(ctx context.Context) ([]*postgres.GetOverdueIssuesRow, error) {
	return r.pg.GetOverdueIssues(ctx)
}

func (r *LibraryRepository) ReturnBook(ctx context.Context, arg postgres.ReturnBookParams) error {
	return r.pg.ReturnBook(ctx, arg)
}

func (r *LibraryRepository) UpdateBookIssue(ctx context.Context, arg postgres.UpdateBookIssueParams) (*postgres.BookIssue, error) {
	return r.pg.UpdateBookIssue(ctx, arg)
}

// ===== Book Ratings Methods =====

func (r *LibraryRepository) CreateBookRating(ctx context.Context, arg postgres.CreateBookRatingParams) (*postgres.BookRating, error) {
	return r.pg.CreateBookRating(ctx, arg)
}

func (r *LibraryRepository) DeleteBookRating(ctx context.Context, ratingID uuid.UUID) error {
	return r.pg.DeleteBookRating(ctx, ratingID)
}

func (r *LibraryRepository) GetBookAverageRating(ctx context.Context, bookID uuid.UUID) (*postgres.GetBookAverageRatingRow, error) {
	return r.pg.GetBookAverageRating(ctx, bookID)
}

func (r *LibraryRepository) GetBookRatingByID(ctx context.Context, ratingID uuid.UUID) (*postgres.GetBookRatingByIDRow, error) {
	return r.pg.GetBookRatingByID(ctx, ratingID)
}

func (r *LibraryRepository) GetBookRatings(ctx context.Context, bookID uuid.UUID) ([]*postgres.GetBookRatingsRow, error) {
	return r.pg.GetBookRatings(ctx, bookID)
}

func (r *LibraryRepository) GetReaderBookRating(ctx context.Context, arg postgres.GetReaderBookRatingParams) (*postgres.BookRating, error) {
	return r.pg.GetReaderBookRating(ctx, arg)
}

func (r *LibraryRepository) GetReaderRatings(ctx context.Context, readerID uuid.UUID) ([]*postgres.GetReaderRatingsRow, error) {
	return r.pg.GetReaderRatings(ctx, readerID)
}

func (r *LibraryRepository) GetTopRatedBooksWithRatings(ctx context.Context, arg postgres.GetTopRatedBooksWithRatingsParams) ([]*postgres.GetTopRatedBooksWithRatingsRow, error) {
	return r.pg.GetTopRatedBooksWithRatings(ctx, arg)
}

func (r *LibraryRepository) UpdateBookRating(ctx context.Context, arg postgres.UpdateBookRatingParams) (*postgres.BookRating, error) {
	return r.pg.UpdateBookRating(ctx, arg)
}

// ===== Books Methods =====

func (r *LibraryRepository) CountBooks(ctx context.Context) (int64, error) {
	return r.pg.CountBooks(ctx)
}

func (r *LibraryRepository) CreateBook(ctx context.Context, arg postgres.CreateBookParams) (*postgres.Book, error) {
	return r.pg.CreateBook(ctx, arg)
}

func (r *LibraryRepository) DeleteBook(ctx context.Context, bookID uuid.UUID) error {
	return r.pg.DeleteBook(ctx, bookID)
}

func (r *LibraryRepository) GetBookByID(ctx context.Context, bookID uuid.UUID) (*postgres.Book, error) {
	return r.pg.GetBookByID(ctx, bookID)
}

func (r *LibraryRepository) GetBookByISBN(ctx context.Context, isbn *string) (*postgres.Book, error) {
	return r.pg.GetBookByISBN(ctx, isbn)
}

func (r *LibraryRepository) GetBookWithDetails(ctx context.Context, bookID uuid.UUID) (*postgres.GetBookWithDetailsRow, error) {
	return r.pg.GetBookWithDetails(ctx, bookID)
}

func (r *LibraryRepository) GetBooksByAuthor(ctx context.Context, authorID uuid.UUID) ([]*postgres.Book, error) {
	return r.pg.GetBooksByAuthor(ctx, authorID)
}

func (r *LibraryRepository) GetBooksWithAuthors(ctx context.Context, arg postgres.GetBooksWithAuthorsParams) ([]*postgres.GetBooksWithAuthorsRow, error) {
	return r.pg.GetBooksWithAuthors(ctx, arg)
}

func (r *LibraryRepository) GetTopRatedBooks(ctx context.Context, arg postgres.GetTopRatedBooksParams) ([]*postgres.GetTopRatedBooksRow, error) {
	return r.pg.GetTopRatedBooks(ctx, arg)
}

func (r *LibraryRepository) ListBooks(ctx context.Context, arg postgres.ListBooksParams) ([]*postgres.Book, error) {
	return r.pg.ListBooks(ctx, arg)
}

func (r *LibraryRepository) SearchBooksByTitle(ctx context.Context, searchQuery string) ([]*postgres.Book, error) {
	return r.pg.SearchBooksByTitle(ctx, searchQuery)
}

func (r *LibraryRepository) UpdateBook(ctx context.Context, arg postgres.UpdateBookParams) (*postgres.Book, error) {
	return r.pg.UpdateBook(ctx, arg)
}

func (r *LibraryRepository) UpdateBookAvailability(ctx context.Context, arg postgres.UpdateBookAvailabilityParams) error {
	return r.pg.UpdateBookAvailability(ctx, arg)
}

// ===== Fines Methods =====

func (r *LibraryRepository) CreateFine(ctx context.Context, arg postgres.CreateFineParams) (*postgres.Fine, error) {
	return r.pg.CreateFine(ctx, arg)
}

func (r *LibraryRepository) DeleteFine(ctx context.Context, fineID uuid.UUID) error {
	return r.pg.DeleteFine(ctx, fineID)
}

func (r *LibraryRepository) GetAllUnpaidFines(ctx context.Context) ([]*postgres.GetAllUnpaidFinesRow, error) {
	return r.pg.GetAllUnpaidFines(ctx)
}

func (r *LibraryRepository) GetFineByID(ctx context.Context, fineID uuid.UUID) (*postgres.GetFineByIDRow, error) {
	return r.pg.GetFineByID(ctx, fineID)
}

func (r *LibraryRepository) GetFineStatistics(ctx context.Context, arg postgres.GetFineStatisticsParams) (*postgres.GetFineStatisticsRow, error) {
	return r.pg.GetFineStatistics(ctx, arg)
}

func (r *LibraryRepository) GetFinesByReader(ctx context.Context, readerID uuid.UUID) ([]*postgres.GetFinesByReaderRow, error) {
	return r.pg.GetFinesByReader(ctx, readerID)
}

func (r *LibraryRepository) GetTotalDebtByReader(ctx context.Context, readerID uuid.UUID) (interface{}, error) {
	return r.pg.GetTotalDebtByReader(ctx, readerID)
}

func (r *LibraryRepository) GetUnpaidFinesByReader(ctx context.Context, readerID uuid.UUID) ([]*postgres.GetUnpaidFinesByReaderRow, error) {
	return r.pg.GetUnpaidFinesByReader(ctx, readerID)
}

func (r *LibraryRepository) PayFine(ctx context.Context, arg postgres.PayFineParams) error {
	return r.pg.PayFine(ctx, arg)
}

func (r *LibraryRepository) UpdateFine(ctx context.Context, arg postgres.UpdateFineParams) (*postgres.Fine, error) {
	return r.pg.UpdateFine(ctx, arg)
}

// ===== Reading Halls Methods =====

func (r *LibraryRepository) CreateReadingHall(ctx context.Context, arg postgres.CreateReadingHallParams) (*postgres.ReadingHall, error) {
	return r.pg.CreateReadingHall(ctx, arg)
}

func (r *LibraryRepository) DeleteReadingHall(ctx context.Context, hallID uuid.UUID) error {
	return r.pg.DeleteReadingHall(ctx, hallID)
}

func (r *LibraryRepository) GetHallStatistics(ctx context.Context) ([]*postgres.GetHallStatisticsRow, error) {
	return r.pg.GetHallStatistics(ctx)
}

func (r *LibraryRepository) GetReadingHallByID(ctx context.Context, hallID uuid.UUID) (*postgres.ReadingHall, error) {
	return r.pg.GetReadingHallByID(ctx, hallID)
}

func (r *LibraryRepository) ListReadingHalls(ctx context.Context) ([]*postgres.ReadingHall, error) {
	return r.pg.ListReadingHalls(ctx)
}

func (r *LibraryRepository) UpdateHallOccupancy(ctx context.Context, arg postgres.UpdateHallOccupancyParams) error {
	return r.pg.UpdateHallOccupancy(ctx, arg)
}

func (r *LibraryRepository) UpdateReadingHall(ctx context.Context, arg postgres.UpdateReadingHallParams) (*postgres.ReadingHall, error) {
	return r.pg.UpdateReadingHall(ctx, arg)
}

// ===== Readers Methods =====

func (r *LibraryRepository) CountReaders(ctx context.Context) (int64, error) {
	return r.pg.CountReaders(ctx)
}

func (r *LibraryRepository) CountReadersByHall(ctx context.Context, hallID *uuid.UUID) (int64, error) {
	return r.pg.CountReadersByHall(ctx, hallID)
}

func (r *LibraryRepository) CreateReader(ctx context.Context, arg postgres.CreateReaderParams) (*postgres.Reader, error) {
	return r.pg.CreateReader(ctx, arg)
}

func (r *LibraryRepository) DeactivateReader(ctx context.Context, readerID uuid.UUID) error {
	return r.pg.DeactivateReader(ctx, readerID)
}

func (r *LibraryRepository) GetReaderByID(ctx context.Context, readerID uuid.UUID) (*postgres.GetReaderByIDRow, error) {
	return r.pg.GetReaderByID(ctx, readerID)
}

func (r *LibraryRepository) GetReaderByTicketNumber(ctx context.Context, ticketNumber string) (*postgres.GetReaderByTicketNumberRow, error) {
	return r.pg.GetReaderByTicketNumber(ctx, ticketNumber)
}

func (r *LibraryRepository) GetReaderByUserID(ctx context.Context, userID *uuid.UUID) (*postgres.GetReaderByUserIDRow, error) {
	return r.pg.GetReaderByUserID(ctx, userID)
}

func (r *LibraryRepository) ListAllReaders(ctx context.Context, arg postgres.ListAllReadersParams) ([]*postgres.ListAllReadersRow, error) {
	return r.pg.ListAllReaders(ctx, arg)
}

func (r *LibraryRepository) ListReadersByHall(ctx context.Context, hallID *uuid.UUID) ([]*postgres.ListReadersByHallRow, error) {
	return r.pg.ListReadersByHall(ctx, hallID)
}

func (r *LibraryRepository) SearchReadersByName(ctx context.Context, searchQuery *string) ([]*postgres.SearchReadersByNameRow, error) {
	return r.pg.SearchReadersByName(ctx, searchQuery)
}

func (r *LibraryRepository) UpdateReader(ctx context.Context, arg postgres.UpdateReaderParams) (*postgres.Reader, error) {
	return r.pg.UpdateReader(ctx, arg)
}

// ===== System Logs Methods =====

func (r *LibraryRepository) CountSystemLogsByAction(ctx context.Context, arg postgres.CountSystemLogsByActionParams) (int64, error) {
	return r.pg.CountSystemLogsByAction(ctx, arg)
}

func (r *LibraryRepository) CreateSystemLog(ctx context.Context, arg postgres.CreateSystemLogParams) error {
	return r.pg.CreateSystemLog(ctx, arg)
}

func (r *LibraryRepository) GetRecentSystemLogs(ctx context.Context, arg postgres.GetRecentSystemLogsParams) ([]*postgres.GetRecentSystemLogsRow, error) {
	return r.pg.GetRecentSystemLogs(ctx, arg)
}

func (r *LibraryRepository) GetSystemLogByID(ctx context.Context, logID uuid.UUID) (*postgres.GetSystemLogByIDRow, error) {
	return r.pg.GetSystemLogByID(ctx, logID)
}

func (r *LibraryRepository) GetSystemLogsByAction(ctx context.Context, arg postgres.GetSystemLogsByActionParams) ([]*postgres.GetSystemLogsByActionRow, error) {
	return r.pg.GetSystemLogsByAction(ctx, arg)
}

func (r *LibraryRepository) GetSystemLogsByEntity(ctx context.Context, arg postgres.GetSystemLogsByEntityParams) ([]*postgres.GetSystemLogsByEntityRow, error) {
	return r.pg.GetSystemLogsByEntity(ctx, arg)
}

func (r *LibraryRepository) GetSystemLogsByUser(ctx context.Context, arg postgres.GetSystemLogsByUserParams) ([]*postgres.GetSystemLogsByUserRow, error) {
	return r.pg.GetSystemLogsByUser(ctx, arg)
}

func (r *LibraryRepository) GetSystemLogsInDateRange(ctx context.Context, arg postgres.GetSystemLogsInDateRangeParams) ([]*postgres.GetSystemLogsInDateRangeRow, error) {
	return r.pg.GetSystemLogsInDateRange(ctx, arg)
}

func (r *LibraryRepository) GetSystemLogsStatistics(ctx context.Context, arg postgres.GetSystemLogsStatisticsParams) (*postgres.GetSystemLogsStatisticsRow, error) {
	return r.pg.GetSystemLogsStatistics(ctx, arg)
}

// ===== Users Methods =====

func (r *LibraryRepository) CountUsers(ctx context.Context) (int64, error) {
	return r.pg.CountUsers(ctx)
}

func (r *LibraryRepository) CountUsersByRole(ctx context.Context, role postgres.UserRole) (int64, error) {
	return r.pg.CountUsersByRole(ctx, role)
}

func (r *LibraryRepository) CreateUser(ctx context.Context, arg postgres.CreateUserParams) (*postgres.User, error) {
	return r.pg.CreateUser(ctx, arg)
}

func (r *LibraryRepository) DeactivateUser(ctx context.Context, userID uuid.UUID) error {
	return r.pg.DeactivateUser(ctx, userID)
}

func (r *LibraryRepository) GetUserByEmail(ctx context.Context, email string) (*postgres.User, error) {
	return r.pg.GetUserByEmail(ctx, email)
}

func (r *LibraryRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*postgres.User, error) {
	return r.pg.GetUserByID(ctx, userID)
}

func (r *LibraryRepository) GetUserByUsername(ctx context.Context, username string) (*postgres.User, error) {
	return r.pg.GetUserByUsername(ctx, username)
}

func (r *LibraryRepository) ListUsers(ctx context.Context, arg postgres.ListUsersParams) ([]*postgres.User, error) {
	return r.pg.ListUsers(ctx, arg)
}

func (r *LibraryRepository) ListUsersByRole(ctx context.Context, role postgres.UserRole) ([]*postgres.User, error) {
	return r.pg.ListUsersByRole(ctx, role)
}

func (r *LibraryRepository) UpdateUser(ctx context.Context, arg postgres.UpdateUserParams) (*postgres.User, error) {
	return r.pg.UpdateUser(ctx, arg)
}

func (r *LibraryRepository) UpdateUserPassword(ctx context.Context, arg postgres.UpdateUserPasswordParams) error {
	return r.pg.UpdateUserPassword(ctx, arg)
}

// ===== Redis Auth Methods =====

func (r *LibraryRepository) CreateSession(ctx context.Context, userID uuid.UUID, role postgres.UserRole, ttl time.Duration) (string, error) {
	return r.rd.CreateSession(ctx, userID, role, ttl)
}

func (r *LibraryRepository) GetSession(ctx context.Context, sessionID string) (redis.Session, error) {
	return r.rd.GetSession(ctx, sessionID)
}

func (r *LibraryRepository) TerminateOtherSessions(ctx context.Context, userID uuid.UUID, currentSessionID string) error {
	return r.rd.TerminateOtherSessions(ctx, userID, currentSessionID)
}

func (r *LibraryRepository) GetUserSessions(ctx context.Context, userID uuid.UUID) ([]redis.Session, error) {
	return r.rd.GetUserSessions(ctx, userID)
}

func (r *LibraryRepository) DeleteSession(ctx context.Context, sessionID string) error {
	return r.rd.DeleteSession(ctx, sessionID)
}

func (r *LibraryRepository) RefreshSession(ctx context.Context, sessionID string, ttl time.Duration) error {
	return r.rd.RefreshSession(ctx, sessionID, ttl)
}

// func (r *LibraryRepository) CreateEmailCode(ctx context.Context, email string, ttl time.Duration) (string, error) {
// 	return r.rd.CreateEmailCode(ctx, email, ttl)
// }

// func (r *LibraryRepository) GetEmailCode(ctx context.Context, email string) (string, error) {
// 	return r.rd.GetEmailCode(ctx, email)
// }

// func (r *LibraryRepository) DeleteEmailCode(ctx context.Context, email string) error {
// 	return r.rd.DeleteEmailCode(ctx, email)
// }

// func (r *LibraryRepository) VerifyEmailCode(ctx context.Context, email string, userProvidedCode string) (bool, error) {
// 	return r.rd.VerifyEmailCode(ctx, email, userProvidedCode)
// }

// func (r *LibraryRepository) MarkEmailAsVerified(ctx context.Context, email string, ttl time.Duration) error {
// 	return r.rd.MarkEmailAsVerified(ctx, email, ttl)
// }

// func (r *LibraryRepository) IsEmailVerified(ctx context.Context, email string) (bool, error) {
// 	return r.rd.IsEmailVerified(ctx, email)
// }
