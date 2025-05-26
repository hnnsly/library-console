package repository

import (
	"context"
	"time"

	"github.com/hnnsly/library-console/internal/repository/postgres"
	"github.com/hnnsly/library-console/internal/repository/redis"
)

// LibraryRepository объединяет все методы для работы с данными библиотеки
type LibraryRepository struct {
	pg *postgres.Queries
	rd *redis.Redis
}

// New создает новый репозиторий библиотеки
func New(pg *postgres.Queries, rd *redis.Redis) *LibraryRepository {
	return &LibraryRepository{
		pg: pg,
		rd: rd,
	}
}

// Books methods
func (r *LibraryRepository) CreateBook(ctx context.Context, params postgres.CreateBookParams) (*postgres.Book, error) {
	return r.pg.CreateBook(ctx, params)
}

func (r *LibraryRepository) GetBookByID(ctx context.Context, bookID int64) (*postgres.GetBookByIDRow, error) {
	return r.pg.GetBookByID(ctx, bookID)
}

func (r *LibraryRepository) GetBookByCode(ctx context.Context, bookCode string) (*postgres.GetBookByCodeRow, error) {
	return r.pg.GetBookByCode(ctx, bookCode)
}

func (r *LibraryRepository) SearchBooks(ctx context.Context, params postgres.SearchBooksParams) ([]*postgres.SearchBooksRow, error) {
	return r.pg.SearchBooks(ctx, params)
}

func (r *LibraryRepository) AdvancedSearchBooks(ctx context.Context, params postgres.AdvancedSearchBooksParams) ([]*postgres.AdvancedSearchBooksRow, error) {
	return r.pg.AdvancedSearchBooks(ctx, params)
}

func (r *LibraryRepository) GetAvailableBooks(ctx context.Context, resultLimit int32) ([]*postgres.GetAvailableBooksRow, error) {
	return r.pg.GetAvailableBooks(ctx, resultLimit)
}

func (r *LibraryRepository) GetPopularBooks(ctx context.Context, limit int32) ([]*postgres.GetPopularBooksRow, error) {
	return r.pg.GetPopularBooks(ctx, limit)
}

func (r *LibraryRepository) GetTopRatedBooks(ctx context.Context, limit int32) ([]*postgres.GetTopRatedBooksRow, error) {
	return r.pg.GetTopRatedBooks(ctx, limit)
}

func (r *LibraryRepository) UpdateBookAvailability(ctx context.Context, params postgres.UpdateBookAvailabilityParams) error {
	return r.pg.UpdateBookAvailability(ctx, params)
}

func (r *LibraryRepository) UpdateBookCopies(ctx context.Context, params postgres.UpdateBookCopiesParams) error {
	return r.pg.UpdateBookCopies(ctx, params)
}

func (r *LibraryRepository) WriteOffBook(ctx context.Context, bookID int64) error {
	return r.pg.WriteOffBook(ctx, bookID)
}

// Readers methods
func (r *LibraryRepository) CreateReader(ctx context.Context, params postgres.CreateReaderParams) (*postgres.Reader, error) {
	return r.pg.CreateReader(ctx, params)
}

func (r *LibraryRepository) GetReaderByID(ctx context.Context, readerID int64) (*postgres.GetReaderByIDRow, error) {
	return r.pg.GetReaderByID(ctx, readerID)
}

func (r *LibraryRepository) GetReaderByTicket(ctx context.Context, ticketNumber string) (*postgres.GetReaderByTicketRow, error) {
	return r.pg.GetReaderByTicket(ctx, ticketNumber)
}

func (r *LibraryRepository) SearchReadersByName(ctx context.Context, params postgres.SearchReadersByNameParams) ([]*postgres.SearchReadersByNameRow, error) {
	return r.pg.SearchReadersByName(ctx, params)
}

func (r *LibraryRepository) GetAllReaders(ctx context.Context, params postgres.GetAllReadersParams) ([]*postgres.GetAllReadersRow, error) {
	return r.pg.GetAllReaders(ctx, params)
}

func (r *LibraryRepository) GetActiveReaders(ctx context.Context, limit int32) ([]*postgres.GetActiveReadersRow, error) {
	return r.pg.GetActiveReaders(ctx, limit)
}

func (r *LibraryRepository) UpdateReader(ctx context.Context, params postgres.UpdateReaderParams) (*postgres.Reader, error) {
	return r.pg.UpdateReader(ctx, params)
}

func (r *LibraryRepository) UpdateReaderStatus(ctx context.Context, params postgres.UpdateReaderStatusParams) error {
	return r.pg.UpdateReaderStatus(ctx, params)
}

func (r *LibraryRepository) UpdateReaderDebt(ctx context.Context, readerID int64) error {
	return r.pg.UpdateReaderDebt(ctx, readerID)
}

func (r *LibraryRepository) GetReaderStatistics(ctx context.Context, readerID int64) (*postgres.GetReaderStatisticsRow, error) {
	return r.pg.GetReaderStatistics(ctx, readerID)
}

func (r *LibraryRepository) GetReaderFavoriteCategories(ctx context.Context, readerID int) ([]*postgres.GetReaderFavoriteCategoriesRow, error) {
	return r.pg.GetReaderFavoriteCategories(ctx, readerID)
}

func (r *LibraryRepository) GetReadersCount(ctx context.Context) (int64, error) {
	return r.pg.GetReadersCount(ctx)
}

// Loans methods
func (r *LibraryRepository) CreateLoan(ctx context.Context, params postgres.CreateLoanParams) (*postgres.LoanHistory, error) {
	return r.pg.CreateLoan(ctx, params)
}

func (r *LibraryRepository) CheckLoanEligibility(ctx context.Context, params postgres.CheckLoanEligibilityParams) (*postgres.CheckLoanEligibilityRow, error) {
	return r.pg.CheckLoanEligibility(ctx, params)
}

func (r *LibraryRepository) ReturnBook(ctx context.Context, params postgres.ReturnBookParams) error {
	return r.pg.ReturnBook(ctx, params)
}

func (r *LibraryRepository) RenewLoan(ctx context.Context, loanId int64) error {
	return r.pg.RenewLoan(ctx, loanId)
}

func (r *LibraryRepository) MarkLoanAsLost(ctx context.Context, loanId int64) error {
	return r.pg.MarkLoanAsLost(ctx, loanId)
}

func (r *LibraryRepository) GetReaderCurrentLoans(ctx context.Context, readerID int) ([]*postgres.GetReaderCurrentLoansRow, error) {
	return r.pg.GetReaderCurrentLoans(ctx, readerID)
}

func (r *LibraryRepository) GetReaderLoanHistory(ctx context.Context, params postgres.GetReaderLoanHistoryParams) ([]*postgres.GetReaderLoanHistoryRow, error) {
	return r.pg.GetReaderLoanHistory(ctx, params)
}

func (r *LibraryRepository) GetOverdueBooks(ctx context.Context, resultLimit int32) ([]*postgres.GetOverdueBooksRow, error) {
	return r.pg.GetOverdueBooks(ctx, resultLimit)
}

func (r *LibraryRepository) GetLoanByID(ctx context.Context, loanID int64) (*postgres.GetLoanByIDRow, error) {
	return r.pg.GetLoanByID(ctx, loanID)
}

func (r *LibraryRepository) GetBooksDueToday(ctx context.Context) ([]*postgres.GetBooksDueTodayRow, error) {
	return r.pg.GetBooksDueToday(ctx)
}

func (r *LibraryRepository) GetActiveLoansByBook(ctx context.Context, bookID int) ([]*postgres.GetActiveLoansByBookRow, error) {
	return r.pg.GetActiveLoansByBook(ctx, bookID)
}

// Renewals methods
func (r *LibraryRepository) CreateRenewal(ctx context.Context, params postgres.CreateRenewalParams) error {
	return r.pg.CreateRenewal(ctx, params)
}

func (r *LibraryRepository) CreateRenewalRecord(ctx context.Context, params postgres.CreateRenewalRecordParams) (*postgres.Renewal, error) {
	return r.pg.CreateRenewalRecord(ctx, params)
}

func (r *LibraryRepository) GetRenewalsForLoan(ctx context.Context, loanHistoryID int) ([]*postgres.GetRenewalsForLoanRow, error) {
	return r.pg.GetRenewalsForLoan(ctx, loanHistoryID)
}

func (r *LibraryRepository) GetRenewalsByDate(ctx context.Context, params postgres.GetRenewalsByDateParams) ([]*postgres.GetRenewalsByDateRow, error) {
	return r.pg.GetRenewalsByDate(ctx, params)
}

func (r *LibraryRepository) GetMostRenewedBooks(ctx context.Context, resultLimit int32) ([]*postgres.GetMostRenewedBooksRow, error) {
	return r.pg.GetMostRenewedBooks(ctx, resultLimit)
}

// Fines methods
func (r *LibraryRepository) CreateFine(ctx context.Context, params postgres.CreateFineParams) (*postgres.Fine, error) {
	return r.pg.CreateFine(ctx, params)
}

func (r *LibraryRepository) CalculateOverdueFine(ctx context.Context, params postgres.CalculateOverdueFineParams) (*postgres.CalculateOverdueFineRow, error) {
	return r.pg.CalculateOverdueFine(ctx, params)
}

func (r *LibraryRepository) PayFine(ctx context.Context, fineID int64) error {
	return r.pg.PayFine(ctx, fineID)
}

func (r *LibraryRepository) WaiveFine(ctx context.Context, fineID int64) error {
	return r.pg.WaiveFine(ctx, fineID)
}

func (r *LibraryRepository) GetReaderFines(ctx context.Context, readerID int) ([]*postgres.GetReaderFinesRow, error) {
	return r.pg.GetReaderFines(ctx, readerID)
}

func (r *LibraryRepository) GetUnpaidFines(ctx context.Context) ([]*postgres.GetUnpaidFinesRow, error) {
	return r.pg.GetUnpaidFines(ctx)
}

func (r *LibraryRepository) GetDebtorReaders(ctx context.Context) ([]*postgres.GetDebtorReadersRow, error) {
	return r.pg.GetDebtorReaders(ctx)
}

// Reservations methods
func (r *LibraryRepository) CreateReservation(ctx context.Context, params postgres.CreateReservationParams) (*postgres.Reservation, error) {
	return r.pg.CreateReservation(ctx, params)
}

func (r *LibraryRepository) FulfillReservation(ctx context.Context, reservationID int64) error {
	return r.pg.FulfillReservation(ctx, reservationID)
}

func (r *LibraryRepository) CancelReservation(ctx context.Context, reservationID int64) error {
	return r.pg.CancelReservation(ctx, reservationID)
}

func (r *LibraryRepository) GetReaderReservations(ctx context.Context, readerID int) ([]*postgres.GetReaderReservationsRow, error) {
	return r.pg.GetReaderReservations(ctx, readerID)
}

func (r *LibraryRepository) GetBookQueue(ctx context.Context, bookID int) ([]*postgres.GetBookQueueRow, error) {
	return r.pg.GetBookQueue(ctx, bookID)
}

func (r *LibraryRepository) GetExpiredReservations(ctx context.Context) ([]*postgres.GetExpiredReservationsRow, error) {
	return r.pg.GetExpiredReservations(ctx)
}

// Halls methods
func (r *LibraryRepository) GetAllHalls(ctx context.Context) ([]*postgres.GetAllHallsRow, error) {
	return r.pg.GetAllHalls(ctx)
}

func (r *LibraryRepository) GetHallByID(ctx context.Context, hallID int64) (*postgres.Hall, error) {
	return r.pg.GetHallByID(ctx, hallID)
}

func (r *LibraryRepository) UpdateHallOccupancy(ctx context.Context, hallID int) error {
	return r.pg.UpdateHallOccupancy(ctx, hallID)
}

func (r *LibraryRepository) GetHallStatistics(ctx context.Context) ([]*postgres.GetHallStatisticsRow, error) {
	return r.pg.GetHallStatistics(ctx)
}

// Categories methods
func (r *LibraryRepository) CreateCategory(ctx context.Context, params postgres.CreateCategoryParams) (*postgres.BookCategory, error) {
	return r.pg.CreateCategory(ctx, params)
}

func (r *LibraryRepository) GetAllCategories(ctx context.Context) ([]*postgres.BookCategory, error) {
	return r.pg.GetAllCategories(ctx)
}

func (r *LibraryRepository) GetCategoryStatistics(ctx context.Context) ([]*postgres.GetCategoryStatisticsRow, error) {
	return r.pg.GetCategoryStatistics(ctx)
}

// Librarians methods
func (r *LibraryRepository) CreateLibrarian(ctx context.Context, params postgres.CreateLibrarianParams) (*postgres.Librarian, error) {
	return r.pg.CreateLibrarian(ctx, params)
}

func (r *LibraryRepository) GetLibrarianByID(ctx context.Context, librarianID int64) (*postgres.Librarian, error) {
	return r.pg.GetLibrarianByID(ctx, librarianID)
}

func (r *LibraryRepository) GetLibrarianByEmployeeID(ctx context.Context, employeeID string) (*postgres.Librarian, error) {
	return r.pg.GetLibrarianByEmployeeID(ctx, employeeID)
}

func (r *LibraryRepository) GetAllLibrarians(ctx context.Context) ([]*postgres.Librarian, error) {
	return r.pg.GetAllLibrarians(ctx)
}

func (r *LibraryRepository) UpdateLibrarian(ctx context.Context, params postgres.UpdateLibrarianParams) (*postgres.Librarian, error) {
	return r.pg.UpdateLibrarian(ctx, params)
}

func (r *LibraryRepository) DeactivateLibrarian(ctx context.Context, librarianID int64) error {
	return r.pg.DeactivateLibrarian(ctx, librarianID)
}

// Operations methods
func (r *LibraryRepository) CreateOperationLog(ctx context.Context, params postgres.CreateOperationLogParams) error {
	return r.pg.CreateOperationLog(ctx, params)
}

func (r *LibraryRepository) GetOperationLogs(ctx context.Context, params postgres.GetOperationLogsParams) ([]*postgres.GetOperationLogsRow, error) {
	return r.pg.GetOperationLogs(ctx, params)
}

func (r *LibraryRepository) GetRecentOperations(ctx context.Context, resultLimit int32) ([]*postgres.GetRecentOperationsRow, error) {
	return r.pg.GetRecentOperations(ctx, resultLimit)
}

// Search methods
func (r *LibraryRepository) GlobalSearch(ctx context.Context, searchTerm string) ([]*postgres.GlobalSearchRow, error) {
	return r.pg.GlobalSearch(ctx, searchTerm)
}

// Statistics and Reports methods
func (r *LibraryRepository) GetLoanStatusStatistics(ctx context.Context, daysBack int) ([]*postgres.GetLoanStatusStatisticsRow, error) {
	date := time.Now().AddDate(0, 0, -daysBack)
	return r.pg.GetLoanStatusStatistics(ctx, date)
}

func (r *LibraryRepository) CreateDailyStatistics(ctx context.Context) error {
	return r.pg.CreateDailyStatistics(ctx)
}

func (r *LibraryRepository) GetMonthlyReport(ctx context.Context) ([]*postgres.GetMonthlyReportRow, error) {
	return r.pg.GetMonthlyReport(ctx)
}

func (r *LibraryRepository) GetYearlyReportByCategory(ctx context.Context) ([]*postgres.GetYearlyReportByCategoryRow, error) {
	return r.pg.GetYearlyReportByCategory(ctx)
}

func (r *LibraryRepository) GetInventoryReport(ctx context.Context) ([]*postgres.GetInventoryReportRow, error) {
	return r.pg.GetInventoryReport(ctx)
}

func (r *LibraryRepository) GetBooksByAuthorInHall(ctx context.Context, params postgres.GetBooksByAuthorInHallParams) (*postgres.GetBooksByAuthorInHallRow, error) {
	return r.pg.GetBooksByAuthorInHall(ctx, params)
}

func (r *LibraryRepository) GetBooksWithSingleCopy(ctx context.Context) ([]*postgres.GetBooksWithSingleCopyRow, error) {
	return r.pg.GetBooksWithSingleCopy(ctx)
}

// Redis cache methods are commented out as in your original file
// func (r *LibraryRepository) CacheBookData(ctx context.Context, bookID int32, data interface{}, ttl time.Duration) error {
// 	return r.rd.Set(ctx, fmt.Sprintf("book:%d", bookID), data, ttl)
// }

// func (r *LibraryRepository) GetCachedBookData(ctx context.Context, bookID int32) (interface{}, error) {
// 	return r.rd.Get(ctx, fmt.Sprintf("book:%d", bookID))
// }

// func (r *LibraryRepository) CacheReaderData(ctx context.Context, readerID int32, data interface{}, ttl time.Duration) error {
// 	return r.rd.Set(ctx, fmt.Sprintf("reader:%d", readerID), data, ttl)
// }

// func (r *LibraryRepository) GetCachedReaderData(ctx context.Context, readerID int32) (interface{}, error) {
// 	return r.rd.Get(ctx, fmt.Sprintf("reader:%d", readerID))
// }
