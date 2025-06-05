import React, { useState } from 'react';
import { motion } from 'framer-motion';
import { Search, Calendar, Check, AlertTriangle, BookOpen } from 'lucide-react';
import { books, members, checkoutRecords } from '../data/mockData';
import { format, isPast, parseISO } from 'date-fns';

const Returns: React.FC = () => {
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedCheckout, setSelectedCheckout] = useState<any | null>(null);
  const [bookCondition, setBookCondition] = useState<'good' | 'damaged'>('good');
  const [fineAmount, setFineAmount] = useState<number>(0);
  const [returnComplete, setReturnComplete] = useState(false);

  // Get active checkouts with book and member details
  const activeCheckouts = checkoutRecords
    .filter(record => record.status === 'active')
    .map(record => {
      const book = books.find(b => b.id === record.bookId);
      const member = members.find(m => m.id === record.memberId);
      
      return {
        ...record,
        book,
        member,
      };
    });

  // Filter checkouts based on search
  const filteredCheckouts = activeCheckouts.filter(checkout => {
    const bookMatches = checkout.book?.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
                       checkout.book?.author.toLowerCase().includes(searchTerm.toLowerCase());
    
    const memberMatches = `${checkout.member?.firstName} ${checkout.member?.lastName}`.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         checkout.member?.email.toLowerCase().includes(searchTerm.toLowerCase());
    
    return bookMatches || memberMatches;
  });

  const handleReturnSubmit = () => {
    if (selectedCheckout) {
      // In a real app, this would update the database
      setReturnComplete(true);
      
      // Reset form after 3 seconds
      setTimeout(() => {
        setSelectedCheckout(null);
        setBookCondition('good');
        setFineAmount(0);
        setReturnComplete(false);
      }, 3000);
    }
  };

  const calculateLateFee = (dueDate: string): number => {
    if (!isPast(parseISO(dueDate))) return 0;
    
    const today = new Date();
    const due = new Date(dueDate);
    const daysLate = Math.floor((today.getTime() - due.getTime()) / (1000 * 60 * 60 * 24));
    
    // $0.50 per day late
    return daysLate * 0.5;
  };

  // Animation variants
  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.1
      }
    }
  };

  const itemVariants = {
    hidden: { y: 20, opacity: 0 },
    visible: {
      y: 0,
      opacity: 1,
      transition: { type: 'spring', stiffness: 300, damping: 24 }
    }
  };

  return (
    <div className="max-w-5xl mx-auto">
      <div className="mb-6">
        <h1 className="text-3xl font-serif font-bold text-gray-900">Return Books</h1>
        <p className="text-gray-600 mt-1">Process book returns and handle any fines or damages.</p>
      </div>

      <motion.div
        variants={containerVariants}
        initial="hidden"
        animate="visible"
        className="bg-white rounded-lg shadow-sm overflow-hidden"
      >
        {returnComplete ? (
          <motion.div 
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            className="p-8 text-center"
          >
            <div className="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <Check size={32} className="text-green-600" />
            </div>
            <h2 className="text-2xl font-bold text-gray-900 mb-2">Return Complete!</h2>
            <p className="text-gray-600 mb-2">
              <span className="font-medium">{selectedCheckout.book?.title}</span> has been returned by <span className="font-medium">{selectedCheckout.member?.firstName} {selectedCheckout.member?.lastName}</span>.
            </p>
            {fineAmount > 0 && (
              <p className="text-orange-700 mb-6">
                A fine of ${fineAmount.toFixed(2)} has been recorded.
              </p>
            )}
            <p className="text-gray-700">
              Book condition: <span className="font-medium">{bookCondition === 'good' ? 'Good' : 'Damaged'}</span>
            </p>
          </motion.div>
        ) : (
          <div>
            <div className="p-6">
              <div className="mb-6">
                <label htmlFor="search" className="label">
                  Search for a book or member
                </label>
                <div className="relative">
                  <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" size={18} />
                  <input
                    id="search"
                    type="text"
                    placeholder="Enter book title, author, or member name..."
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                    className="input pl-10 py-2"
                  />
                </div>
              </div>
              
              <motion.div variants={itemVariants} className="mb-6">
                <h2 className="text-xl font-semibold mb-4 flex items-center">
                  <BookOpen size={20} className="mr-2 text-primary-500" /> 
                  Active Checkouts
                </h2>
                
                {filteredCheckouts.length > 0 ? (
                  <div className="overflow-x-auto border rounded-lg">
                    <table className="min-w-full divide-y divide-gray-200">
                      <thead className="bg-gray-50">
                        <tr>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Book</th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Member</th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Checkout Date</th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Due Date</th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                          <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Action</th>
                        </tr>
                      </thead>
                      <tbody className="bg-white divide-y divide-gray-200">
                        {filteredCheckouts.map((checkout) => {
                          const isOverdue = isPast(new Date(checkout.dueDate));
                          
                          return (
                            <tr 
                              key={checkout.id} 
                              className={`hover:bg-gray-50 cursor-pointer ${
                                selectedCheckout?.id === checkout.id ? 'bg-primary-50' : ''
                              }`}
                              onClick={() => {
                                setSelectedCheckout(checkout);
                                const lateFee = calculateLateFee(checkout.dueDate);
                                setFineAmount(lateFee);
                              }}
                            >
                              <td className="px-6 py-4 whitespace-nowrap">
                                <div className="flex items-center">
                                  <div className="h-10 w-7 bg-gray-200 rounded overflow-hidden mr-3 flex-shrink-0">
                                    {checkout.book?.coverImage && (
                                      <img 
                                        src={checkout.book.coverImage} 
                                        alt={checkout.book.title} 
                                        className="h-full w-full object-cover"
                                      />
                                    )}
                                  </div>
                                  <div>
                                    <div className="font-medium text-gray-900">{checkout.book?.title}</div>
                                    <div className="text-sm text-gray-500">{checkout.book?.author}</div>
                                  </div>
                                </div>
                              </td>
                              <td className="px-6 py-4 whitespace-nowrap">
                                <div className="flex items-center">
                                  <div className="text-sm font-medium text-gray-900">
                                    {checkout.member?.firstName} {checkout.member?.lastName}
                                  </div>
                                </div>
                              </td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                {format(new Date(checkout.checkoutDate), 'MMM d, yyyy')}
                              </td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                {format(new Date(checkout.dueDate), 'MMM d, yyyy')}
                              </td>
                              <td className="px-6 py-4 whitespace-nowrap">
                                {isOverdue ? (
                                  <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800">
                                    <AlertTriangle size={12} className="mr-1" /> Overdue
                                  </span>
                                ) : (
                                  <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                                    <Check size={12} className="mr-1" /> On time
                                  </span>
                                )}
                              </td>
                              <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                                <button 
                                  className="text-primary-600 hover:text-primary-900"
                                  onClick={(e) => {
                                    e.stopPropagation();
                                    setSelectedCheckout(checkout);
                                    const lateFee = calculateLateFee(checkout.dueDate);
                                    setFineAmount(lateFee);
                                  }}
                                >
                                  Process Return
                                </button>
                              </td>
                            </tr>
                          );
                        })}
                      </tbody>
                    </table>
                  </div>
                ) : (
                  <div className="text-center py-8 border rounded-lg">
                    <p className="text-gray-500">
                      {searchTerm ? 'No matching checkouts found.' : 'No active checkouts to display.'}
                    </p>
                  </div>
                )}
              </motion.div>
              
              {selectedCheckout && (
                <motion.div 
                  initial={{ opacity: 0, height: 0 }}
                  animate={{ opacity: 1, height: 'auto' }}
                  transition={{ duration: 0.3 }}
                  className="border rounded-lg p-6 bg-gray-50"
                >
                  <h2 className="text-xl font-semibold mb-4">Process Return</h2>
                  
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div>
                      <h3 className="font-medium text-gray-800 mb-3">Book Details</h3>
                      <div className="flex items-start">
                        <div className="h-24 w-16 bg-gray-200 rounded overflow-hidden mr-4 flex-shrink-0">
                          {selectedCheckout.book?.coverImage && (
                            <img 
                              src={selectedCheckout.book.coverImage} 
                              alt={selectedCheckout.book.title} 
                              className="h-full w-full object-cover"
                            />
                          )}
                        </div>
                        <div>
                          <p className="font-medium text-gray-900">{selectedCheckout.book?.title}</p>
                          <p className="text-sm text-gray-600 mb-1">{selectedCheckout.book?.author}</p>
                          <p className="text-sm text-gray-600">ISBN: {selectedCheckout.book?.isbn}</p>
                          <p className="text-sm text-gray-600 mt-2">
                            Due date: {format(new Date(selectedCheckout.dueDate), 'MMMM d, yyyy')}
                          </p>
                          {isPast(new Date(selectedCheckout.dueDate)) && (
                            <p className="text-sm text-red-600 font-medium mt-1">
                              This book is overdue!
                            </p>
                          )}
                        </div>
                      </div>
                    </div>
                    
                    <div>
                      <h3 className="font-medium text-gray-800 mb-3">Return Details</h3>
                      
                      <div className="mb-4">
                        <label className="label">Book Condition</label>
                        <div className="flex space-x-4">
                          <label className="flex items-center">
                            <input
                              type="radio"
                              name="condition"
                              value="good"
                              checked={bookCondition === 'good'}
                              onChange={() => setBookCondition('good')}
                              className="mr-2"
                            />
                            <span>Good</span>
                          </label>
                          <label className="flex items-center">
                            <input
                              type="radio"
                              name="condition"
                              value="damaged"
                              checked={bookCondition === 'damaged'}
                              onChange={() => {
                                setBookCondition('damaged');
                                // Add damage fee if condition is damaged
                                const lateFee = calculateLateFee(selectedCheckout.dueDate);
                                setFineAmount(lateFee + 10); // $10 damage fee
                              }}
                              className="mr-2"
                            />
                            <span>Damaged</span>
                          </label>
                        </div>
                      </div>
                      
                      <div className="mb-4">
                        <label htmlFor="fineAmount" className="label">Fine Amount ($)</label>
                        <input
                          id="fineAmount"
                          type="number"
                          step="0.01"
                          min="0"
                          value={fineAmount}
                          onChange={(e) => setFineAmount(parseFloat(e.target.value))}
                          className="input py-2"
                        />
                        {isPast(new Date(selectedCheckout.dueDate)) && (
                          <p className="text-xs text-gray-500 mt-1">
                            Includes late fee of ${calculateLateFee(selectedCheckout.dueDate).toFixed(2)}
                          </p>
                        )}
                        {bookCondition === 'damaged' && (
                          <p className="text-xs text-gray-500 mt-1">
                            Includes damage fee of $10.00
                          </p>
                        )}
                      </div>
                    </div>
                  </div>
                </motion.div>
              )}
            </div>
            
            <div className="bg-gray-50 p-6 flex justify-end border-t">
              <button
                className={`btn flex items-center ${
                  selectedCheckout 
                    ? 'bg-primary-500 text-white hover:bg-primary-600' 
                    : 'bg-gray-200 text-gray-500 cursor-not-allowed'
                }`}
                disabled={!selectedCheckout}
                onClick={handleReturnSubmit}
              >
                Complete Return <Check size={18} className="ml-2" />
              </button>
            </div>
          </div>
        )}
      </motion.div>
    </div>
  );
};

export default Returns;