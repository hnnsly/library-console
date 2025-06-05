import React from "react";
import { useParams, Link, useNavigate } from "react-router-dom";
import { motion } from "framer-motion";
import {
  ArrowLeft,
  Edit,
  Trash2,
  Mail,
  Phone,
  MapPin,
  Calendar,
  AlertTriangle,
  Clock,
  Check,
} from "lucide-react";
import { members, books } from "../data/mockData";

const MemberDetails: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const member = members.find((member) => member.id === id);

  if (!member) {
    return (
      <div className="max-w-4xl mx-auto text-center py-12">
        <h1 className="text-2xl font-bold text-gray-900 mb-4">
          Member Not Found
        </h1>
        <p className="text-gray-600 mb-6">
          The member you're looking for doesn't exist or has been removed.
        </p>
        <button
          onClick={() => navigate("/members")}
          className="btn btn-primary"
        >
          Back to Members
        </button>
      </div>
    );
  }

  // Get the book details for borrowed books
  const borrowedBooksDetails = member.borrowedBooks.map((borrowedBook) => {
    const book = books.find((b) => b.id === borrowedBook.bookId);
    return {
      ...borrowedBook,
      bookDetails: book,
    };
  });

  // Get borrowing history with book details
  const borrowingHistoryWithDetails = member.borrowingHistory.map(
    (historyItem) => {
      const book = books.find((b) => b.id === historyItem.bookId);
      return {
        ...historyItem,
        bookDetails: book,
      };
    },
  );

  // Sort borrowing history by return date (most recent first)
  const sortedBorrowingHistory = [...borrowingHistoryWithDetails].sort(
    (a, b) =>
      new Date(b.returnDate).getTime() - new Date(a.returnDate).getTime(),
  );

  // Animation variants
  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.1,
      },
    },
  };

  const itemVariants = {
    hidden: { y: 20, opacity: 0 },
    visible: {
      y: 0,
      opacity: 1,
      transition: { type: "spring", stiffness: 300, damping: 24 },
    },
  };

  return (
    <div className="max-w-7xl mx-auto">
      <div className="mb-6">
        <button
          onClick={() => navigate("/members")}
          className="flex items-center text-primary-600 hover:underline"
        >
          <ArrowLeft size={16} className="mr-1" /> Back to Members
        </button>
      </div>

      <motion.div
        variants={containerVariants}
        initial="hidden"
        animate="visible"
      >
        <motion.div
          variants={itemVariants}
          className="bg-white rounded-lg shadow-sm overflow-hidden mb-6"
        >
          <div
            className={`h-2 ${
              member.membershipStatus === "active"
                ? "bg-green-500"
                : member.membershipStatus === "expired"
                  ? "bg-orange-500"
                  : "bg-red-500"
            }`}
          ></div>

          <div className="p-6">
            <div className="flex flex-col md:flex-row md:items-center justify-between mb-6">
              <div className="flex items-center">
                <div className="w-16 h-16 rounded-full bg-primary-500 text-white flex items-center justify-center text-2xl font-medium">
                  {member.firstName.charAt(0)}
                  {member.lastName.charAt(0)}
                </div>
                <div className="ml-4">
                  <div className="flex items-center">
                    <h1 className="text-3xl font-serif font-bold">
                      {member.firstName} {member.lastName}
                    </h1>
                    {member.membershipStatus === "active" ? (
                      <span className="ml-3 inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                        Active
                      </span>
                    ) : member.membershipStatus === "expired" ? (
                      <span className="ml-3 inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-orange-100 text-orange-800">
                        Expired
                      </span>
                    ) : (
                      <span className="ml-3 inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800">
                        Suspended
                      </span>
                    )}
                  </div>
                  <p className="text-gray-600">Member ID: {member.id}</p>
                </div>
              </div>
              <div className="mt-4 md:mt-0 flex space-x-3">
                <Link
                  to={`/members/${member.id}/edit`}
                  className="btn btn-secondary flex items-center"
                >
                  <Edit size={16} className="mr-2" /> Edit
                </Link>
                <button className="btn bg-red-500 text-white hover:bg-red-600 flex items-center">
                  <Trash2 size={16} className="mr-2" /> Delete
                </button>
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <div className="flex items-start">
                <div className="p-2 rounded-full bg-gray-100 text-gray-600">
                  <Mail size={18} />
                </div>
                <div className="ml-3">
                  <h3 className="text-sm text-gray-500">Email</h3>
                  <p className="font-medium">{member.email}</p>
                </div>
              </div>

              <div className="flex items-start">
                <div className="p-2 rounded-full bg-gray-100 text-gray-600">
                  <Phone size={18} />
                </div>
                <div className="ml-3">
                  <h3 className="text-sm text-gray-500">Phone</h3>
                  <p className="font-medium">{member.phone}</p>
                </div>
              </div>

              <div className="flex items-start">
                <div className="p-2 rounded-full bg-gray-100 text-gray-600">
                  <Calendar size={18} />
                </div>
                <div className="ml-3">
                  <h3 className="text-sm text-gray-500">Join Date</h3>
                  <p className="font-medium">
                    {new Date(member.joinDate).toLocaleDateString()}
                  </p>
                </div>
              </div>
            </div>

            <div className="mt-6">
              <div className="flex items-start">
                <div className="p-2 rounded-full bg-gray-100 text-gray-600">
                  <MapPin size={18} />
                </div>
                <div className="ml-3">
                  <h3 className="text-sm text-gray-500">Address</h3>
                  <p className="font-medium">{member.address}</p>
                </div>
              </div>
            </div>

            {member.fines > 0 && (
              <div className="mt-6 bg-red-50 border border-red-100 rounded-lg p-4 flex items-center">
                <AlertTriangle size={20} className="text-red-500 mr-3" />
                <div>
                  <h3 className="font-medium text-red-800">
                    Outstanding Fines
                  </h3>
                  <p className="text-red-700">
                    This member has ${member.fines.toFixed(2)} in unpaid fines.
                  </p>
                </div>
                <button className="ml-auto btn bg-white text-red-600 border border-red-200 hover:bg-red-50">
                  Collect Payment
                </button>
              </div>
            )}
          </div>
        </motion.div>

        {/* Currently Borrowed Books */}
        <motion.div
          variants={itemVariants}
          className="bg-white rounded-lg shadow-sm overflow-hidden mb-6"
        >
          <div className="p-4 bg-primary-500 text-white font-medium">
            Currently Borrowed Books
          </div>
          <div className="p-6">
            {borrowedBooksDetails.length > 0 ? (
              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead>
                    <tr>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Book
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Borrowed Date
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Due Date
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Status
                      </th>
                      <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Actions
                      </th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-200">
                    {borrowedBooksDetails.map((item) => {
                      const isOverdue = new Date(item.dueDate) < new Date();

                      return (
                        <tr key={item.bookId}>
                          <td className="px-4 py-4 whitespace-nowrap">
                            <Link
                              to={`/books/${item.bookId}`}
                              className="text-primary-600 hover:underline"
                            >
                              {item.bookDetails?.title}
                            </Link>
                          </td>
                          <td className="px-4 py-4 whitespace-nowrap text-gray-700">
                            {new Date(item.borrowDate).toLocaleDateString()}
                          </td>
                          <td className="px-4 py-4 whitespace-nowrap text-gray-700">
                            {new Date(item.dueDate).toLocaleDateString()}
                          </td>
                          <td className="px-4 py-4 whitespace-nowrap">
                            {isOverdue ? (
                              <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800">
                                Overdue
                              </span>
                            ) : (
                              <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                                On Time
                              </span>
                            )}
                          </td>
                          <td className="px-4 py-4 whitespace-nowrap text-right">
                            <button className="text-primary-600 hover:text-primary-900 mr-3">
                              Renew
                            </button>
                            <button className="text-green-600 hover:text-green-900">
                              Return
                            </button>
                          </td>
                        </tr>
                      );
                    })}
                  </tbody>
                </table>
              </div>
            ) : (
              <p className="text-gray-500 text-center py-6">
                This member has no books currently checked out.
              </p>
            )}
          </div>
        </motion.div>

        {/* Borrowing History */}
        <motion.div
          variants={itemVariants}
          className="bg-white rounded-lg shadow-sm overflow-hidden"
        >
          <div className="p-4 bg-accent-500 text-white font-medium">
            Borrowing History
          </div>
          <div className="p-6">
            {sortedBorrowingHistory.length > 0 ? (
              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead>
                    <tr>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Book
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Borrowed
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Returned
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Status
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Condition
                      </th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-200">
                    {sortedBorrowingHistory.map((item, index) => (
                      <tr key={index}>
                        <td className="px-4 py-4 whitespace-nowrap">
                          <Link
                            to={`/books/${item.bookId}`}
                            className="text-primary-600 hover:underline"
                          >
                            {item.bookDetails?.title}
                          </Link>
                        </td>
                        <td className="px-4 py-4 whitespace-nowrap text-gray-700">
                          {new Date(item.borrowDate).toLocaleDateString()}
                        </td>
                        <td className="px-4 py-4 whitespace-nowrap text-gray-700">
                          {new Date(item.returnDate).toLocaleDateString()}
                        </td>
                        <td className="px-4 py-4 whitespace-nowrap">
                          {item.wasLate ? (
                            <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800">
                              <Clock size={12} className="mr-1" /> Late
                            </span>
                          ) : (
                            <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                              <Check size={12} className="mr-1" /> On Time
                            </span>
                          )}
                        </td>
                        <td className="px-4 py-4 whitespace-nowrap">
                          {item.condition === "good" ? (
                            <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                              Good
                            </span>
                          ) : (
                            <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-orange-100 text-orange-800">
                              Damaged
                            </span>
                          )}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            ) : (
              <p className="text-gray-500 text-center py-6">
                This member has no borrowing history.
              </p>
            )}
          </div>
        </motion.div>
      </motion.div>
    </div>
  );
};

export default MemberDetails;
