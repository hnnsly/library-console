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

func (r *LibraryRepository) GetBookByID(ctx context.Context, bookID int32) (*postgres.GetBookByIDRow, error) {
	return r.pg.GetBookByID(ctx, bookID)
}

func (r *LibraryRepository) GetBookByCode(ctx context.Context, bookCode string) (*postgres.GetBookByCodeRow, error) {
	return r.pg.GetBookByCode(ctx, bookCode)
}

func (r *LibraryRepository) SearchBooks(ctx context.Context, params postgres.SearchBooksParams) ([]*postgres.SearchBooksRow, error) {
	return r.pg.SearchBooks(ctx, params)
}

func (r *LibraryRepository) GetAvailableBooks(ctx context.Context, resultLimit int32) ([]*postgres.GetAvailableBooksRow, error) {
	return r.pg.GetAvailableBooks(ctx, resultLimit)
}

func (r *LibraryRepository) GetPopularBooks(ctx context.Context, limit int32) ([]*postgres.GetPopularBooksRow, error) {
	return r.pg.GetPopularBooks(ctx, limit)
}

func (r *LibraryRepository) UpdateBookAvailability(ctx context.Context, params postgres.UpdateBookAvailabilityParams) error {
	return r.pg.UpdateBookAvailability(ctx, params)
}

// Readers methods
func (r *LibraryRepository) CreateReader(ctx context.Context, params postgres.CreateReaderParams) (*postgres.Reader, error) {
	return r.pg.CreateReader(ctx, params)
}

func (r *LibraryRepository) GetReaderByID(ctx context.Context, readerID int32) (*postgres.GetReaderByIDRow, error) {
	return r.pg.GetReaderByID(ctx, readerID)
}

func (r *LibraryRepository) GetReaderByTicket(ctx context.Context, ticketNumber string) (*postgres.GetReaderByTicketRow, error) {
	return r.pg.GetReaderByTicket(ctx, ticketNumber)
}

func (r *LibraryRepository) SearchReadersByName(ctx context.Context, params postgres.SearchReadersByNameParams) ([]*postgres.SearchReadersByNameRow, error) {
	return r.pg.SearchReadersByName(ctx, params)
}

func (r *LibraryRepository) UpdateReader(ctx context.Context, params postgres.UpdateReaderParams) (*postgres.Reader, error) {
	return r.pg.UpdateReader(ctx, params)
}

func (r *LibraryRepository) UpdateReaderStatus(ctx context.Context, params postgres.UpdateReaderStatusParams) error {
	return r.pg.UpdateReaderStatus(ctx, params)
}

func (r *LibraryRepository) GetReaderStatistics(ctx context.Context, readerID int32) (*postgres.GetReaderStatisticsRow, error) {
	return r.pg.GetReaderStatistics(ctx, readerID)
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

func (r *LibraryRepository) RenewLoan(ctx context.Context, loanId int32) error {
	return r.pg.RenewLoan(ctx, loanId)
}

func (r *LibraryRepository) GetReaderCurrentLoans(ctx context.Context, readerID int32) ([]*postgres.GetReaderCurrentLoansRow, error) {
	return r.pg.GetReaderCurrentLoans(ctx, readerID)
}

func (r *LibraryRepository) GetReaderLoanHistory(ctx context.Context, params postgres.GetReaderLoanHistoryParams) ([]*postgres.GetReaderLoanHistoryRow, error) {
	return r.pg.GetReaderLoanHistory(ctx, params)
}

func (r *LibraryRepository) GetOverdueBooks(ctx context.Context, resultLimit int32) ([]*postgres.GetOverdueBooksRow, error) {
	return r.pg.GetOverdueBooks(ctx, resultLimit)
}

func (r *LibraryRepository) GetLoanByID(ctx context.Context, loanID int32) (*postgres.GetLoanByIDRow, error) {
	return r.pg.GetLoanByID(ctx, loanID)
}

// Fines methods
func (r *LibraryRepository) CreateFine(ctx context.Context, params postgres.CreateFineParams) (*postgres.Fine, error) {
	return r.pg.CreateFine(ctx, params)
}

func (r *LibraryRepository) PayFine(ctx context.Context, fineID int32) error {
	return r.pg.PayFine(ctx, fineID)
}

func (r *LibraryRepository) WaiveFine(ctx context.Context, fineID int32) error {
	return r.pg.WaiveFine(ctx, fineID)
}

func (r *LibraryRepository) GetReaderFines(ctx context.Context, readerID int32) ([]*postgres.GetReaderFinesRow, error) {
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

func (r *LibraryRepository) FulfillReservation(ctx context.Context, reservationID int32) error {
	return r.pg.FulfillReservation(ctx, reservationID)
}

func (r *LibraryRepository) CancelReservation(ctx context.Context, reservationID int32) error {
	return r.pg.CancelReservation(ctx, reservationID)
}

func (r *LibraryRepository) GetReaderReservations(ctx context.Context, readerID int32) ([]*postgres.GetReaderReservationsRow, error) {
	return r.pg.GetReaderReservations(ctx, readerID)
}

func (r *LibraryRepository) GetBookQueue(ctx context.Context, bookID int32) ([]*postgres.GetBookQueueRow, error) {
	return r.pg.GetBookQueue(ctx, bookID)
}

// Halls methods
func (r *LibraryRepository) GetAllHalls(ctx context.Context) ([]*postgres.GetAllHallsRow, error) {
	return r.pg.GetAllHalls(ctx)
}

func (r *LibraryRepository) GetHallByID(ctx context.Context, hallID int32) (*postgres.Hall, error) {
	return r.pg.GetHallByID(ctx, hallID)
}

func (r *LibraryRepository) UpdateHallOccupancy(ctx context.Context, hallID int32) error {
	return r.pg.UpdateHallOccupancy(ctx, hallID)
}

// Categories methods
func (r *LibraryRepository) CreateCategory(ctx context.Context, params postgres.CreateCategoryParams) (*postgres.BookCategory, error) {
	return r.pg.CreateCategory(ctx, params)
}

func (r *LibraryRepository) GetAllCategories(ctx context.Context) ([]*postgres.BookCategory, error) {
	return r.pg.GetAllCategories(ctx)
}

// Librarians methods
func (r *LibraryRepository) CreateLibrarian(ctx context.Context, params postgres.CreateLibrarianParams) (*postgres.Librarian, error) {
	return r.pg.CreateLibrarian(ctx, params)
}

func (r *LibraryRepository) GetLibrarianByID(ctx context.Context, librarianID int32) (*postgres.Librarian, error) {
	return r.pg.GetLibrarianByID(ctx, librarianID)
}

func (r *LibraryRepository) GetAllLibrarians(ctx context.Context) ([]*postgres.Librarian, error) {
	return r.pg.GetAllLibrarians(ctx)
}

func (r *LibraryRepository) UpdateLibrarian(ctx context.Context, params postgres.UpdateLibrarianParams) (*postgres.Librarian, error) {
	return r.pg.UpdateLibrarian(ctx, params)
}

// Operations methods
func (r *LibraryRepository) CreateOperationLog(ctx context.Context, params postgres.CreateOperationLogParams) error {
	return r.pg.CreateOperationLog(ctx, params)
}

func (r *LibraryRepository) GetOperationLogs(ctx context.Context, params postgres.GetOperationLogsParams) ([]*postgres.GetOperationLogsRow, error) {
	return r.pg.GetOperationLogs(ctx, params)
}

// Search methods
func (r *LibraryRepository) GlobalSearch(ctx context.Context, searchTerm string) ([]*postgres.GlobalSearchRow, error) {
	return r.pg.GlobalSearch(ctx, searchTerm)
}

// Statistics methods
func (r *LibraryRepository) GetLoanStatusStatistics(ctx context.Context, daysBack int) ([]*postgres.GetLoanStatusStatisticsRow, error) {
	date := time.Now().AddDate(0, 0, -daysBack)
	return r.pg.GetLoanStatusStatistics(ctx, date)
}

func (r *LibraryRepository) CreateDailyStatistics(ctx context.Context) error {
	return r.pg.CreateDailyStatistics(ctx)
}

// Redis cache methods
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
