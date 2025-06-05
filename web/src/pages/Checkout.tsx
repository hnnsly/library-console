import React, { useState } from 'react';
import { motion } from 'framer-motion';
import { Search, BookOpen, User, ArrowRight, Calendar, Check } from 'lucide-react';
import { books, members } from '../data/mockData';
import { addDays, format } from 'date-fns';

const Checkout: React.FC = () => {
  const [bookSearchTerm, setBookSearchTerm] = useState('');
  const [memberSearchTerm, setMemberSearchTerm] = useState('');
  const [selectedBook, setSelectedBook] = useState<any | null>(null);
  const [selectedMember, setSelectedMember] = useState<any | null>(null);
  const [dueDate, setDueDate] = useState<string>(
    format(addDays(new Date(), 14), 'yyyy-MM-dd')
  );
  const [checkoutComplete, setCheckoutComplete] = useState(false);

  // Filter books based on search term
  const filteredBooks = books
    .filter(book => 
      book.status === 'available' && book.availableCopies > 0 &&
      (book.title.toLowerCase().includes(bookSearchTerm.toLowerCase()) || 
       book.author.toLowerCase().includes(bookSearchTerm.toLowerCase()) ||
       book.isbn.includes(bookSearchTerm))
    )
    .slice(0, 5);

  // Filter members based on search term
  const filteredMembers = members
    .filter(member => 
      member.membershipStatus === 'active' &&
      (`${member.firstName} ${member.lastName}`.toLowerCase().includes(memberSearchTerm.toLowerCase()) ||
       member.email.toLowerCase().includes(memberSearchTerm.toLowerCase()) ||
       member.phone.includes(memberSearchTerm))
    )
    .slice(0, 5);

  const handleCheckout = () => {
    if (selectedBook && selectedMember) {
      // In a real app, this would update the database
      setCheckoutComplete(true);
      
      // Reset form after 3 seconds
      setTimeout(() => {
        setSelectedBook(null);
        setSelectedMember(null);
        setDueDate(format(addDays(new Date(), 14), 'yyyy-MM-dd'));
        setBookSearchTerm('');
        setMemberSearchTerm('');
        setCheckoutComplete(false);
      }, 3000);
    }
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
        <h1 className="text-3xl font-serif font-bold text-gray-900">Checkout Books</h1>
        <p className="text-gray-600 mt-1">Issue books to library members.</p>
      </div>

      <motion.div
        variants={containerVariants}
        initial="hidden"
        animate="visible"
        className="bg-white rounded-lg shadow-sm overflow-hidden"
      >
        {checkoutComplete ? (
          <motion.div 
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            className="p-8 text-center"
          >
            <div className="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <Check size={32} className="text-green-600" />
            </div>
            <h2 className="text-2xl font-bold text-gray-900 mb-2">Checkout Complete!</h2>
            <p className="text-gray-600 mb-6">
              <span className="font-medium">{selectedBook.title}</span> has been checked out to <span className="font-medium">{selectedMember.firstName} {selectedMember.lastName}</span>.
            </p>
            <p className="text-gray-700">
              Due date: <span className="font-medium">{format(new Date(dueDate), 'MMMM d, yyyy')}</span>
            </p>
          </motion.div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 divide-y md:divide-y-0 md:divide-x divide-gray-200">
            {/* Left Column: Book Selection */}
            <motion.div variants={itemVariants} className="p-6">
              <h2 className="text-xl font-semibold mb-4 flex items-center">
                <BookOpen size={20} className="mr-2 text-primary-500" /> 
                Select Book
              </h2>
              
              <div className="mb-4">
                <label htmlFor="bookSearch" className="label">
                  Search for a book by title, author, or ISBN
                </label>
                <div className="relative">
                  <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" size={18} />
                  <input
                    id="bookSearch"
                    type="text"
                    placeholder="Start typing..."
                    value={bookSearchTerm}
                    onChange={(e) => setBookSearchTerm(e.target.value)}
                    className="input pl-10 py-2"
                  />
                </div>
              </div>
              
              {bookSearchTerm && (
                <div className="mb-4 border rounded-lg overflow-hidden">
                  {filteredBooks.length > 0 ? (
                    <ul className="divide-y divide-gray-200">
                      {filteredBooks.map((book) => (
                        <li 
                          key={book.id}
                          className={`p-3 cursor-pointer hover:bg-gray-50 transition-colors ${
                            selectedBook?.id === book.id ? 'bg-primary-50' : ''
                          }`}
                          onClick={() => setSelectedBook(book)}
                        >
                          <div className="flex items-start">
                            <div className="h-12 w-8 bg-gray-200 rounded overflow-hidden mr-3 flex-shrink-0">
                              {book.coverImage && (
                                <img 
                                  src={book.coverImage} 
                                  alt={book.title} 
                                  className="h-full w-full object-cover"
                                />
                              )}
                            </div>
                            <div>
                              <p className="font-medium text-gray-900">{book.title}</p>
                              <p className="text-sm text-gray-600">{book.author}</p>
                              <div className="flex items-center mt-1">
                                <span className="text-xs bg-green-100 text-green-800 px-2 py-0.5 rounded-full">
                                  {book.availableCopies} available
                                </span>
                              </div>
                            </div>
                          </div>
                        </li>
                      ))}
                    </ul>
                  ) : (
                    <div className="p-4 text-center text-gray-500">
                      No available books found. Try a different search term.
                    </div>
                  )}
                </div>
              )}
              
              {selectedBook && (
                <div className="mt-6 p-4 border border-primary-100 rounded-lg bg-primary-50">
                  <h3 className="font-medium text-primary-900 mb-2">Selected Book</h3>
                  <div className="flex items-start">
                    <div className="h-24 w-16 bg-gray-200 rounded overflow-hidden mr-4 flex-shrink-0">
                      {selectedBook.coverImage && (
                        <img 
                          src={selectedBook.coverImage} 
                          alt={selectedBook.title} 
                          className="h-full w-full object-cover"
                        />
                      )}
                    </div>
                    <div>
                      <p className="font-medium text-gray-900">{selectedBook.title}</p>
                      <p className="text-sm text-gray-600 mb-1">{selectedBook.author}</p>
                      <p className="text-sm text-gray-600">ISBN: {selectedBook.isbn}</p>
                      <p className="text-sm text-gray-600">
                        Location: {selectedBook.location}
                      </p>
                      <button 
                        onClick={() => setSelectedBook(null)}
                        className="text-sm text-red-600 hover:underline mt-2"
                      >
                        Remove
                      </button>
                    </div>
                  </div>
                </div>
              )}
            </motion.div>
            
            {/* Right Column: Member Selection */}
            <motion.div variants={itemVariants} className="p-6">
              <h2 className="text-xl font-semibold mb-4 flex items-center">
                <User size={20} className="mr-2 text-primary-500" /> 
                Select Member
              </h2>
              
              <div className="mb-4">
                <label htmlFor="memberSearch" className="label">
                  Search for a member by name, email, or phone
                </label>
                <div className="relative">
                  <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" size={18} />
                  <input
                    id="memberSearch"
                    type="text"
                    placeholder="Start typing..."
                    value={memberSearchTerm}
                    onChange={(e) => setMemberSearchTerm(e.target.value)}
                    className="input pl-10 py-2"
                  />
                </div>
              </div>
              
              {memberSearchTerm && (
                <div className="mb-4 border rounded-lg overflow-hidden">
                  {filteredMembers.length > 0 ? (
                    <ul className="divide-y divide-gray-200">
                      {filteredMembers.map((member) => (
                        <li 
                          key={member.id}
                          className={`p-3 cursor-pointer hover:bg-gray-50 transition-colors ${
                            selectedMember?.id === member.id ? 'bg-primary-50' : ''
                          }`}
                          onClick={() => setSelectedMember(member)}
                        >
                          <div className="flex items-center">
                            <div className="w-10 h-10 rounded-full bg-primary-500 text-white flex items-center justify-center mr-3">
                              {member.firstName.charAt(0)}{member.lastName.charAt(0)}
                            </div>
                            <div>
                              <p className="font-medium text-gray-900">
                                {member.firstName} {member.lastName}
                              </p>
                              <p className="text-sm text-gray-600">{member.email}</p>
                            </div>
                          </div>
                        </li>
                      ))}
                    </ul>
                  ) : (
                    <div className="p-4 text-center text-gray-500">
                      No active members found. Try a different search term.
                    </div>
                  )}
                </div>
              )}
              
              {selectedMember && (
                <div className="mt-6 p-4 border border-primary-100 rounded-lg bg-primary-50">
                  <h3 className="font-medium text-primary-900 mb-2">Selected Member</h3>
                  <div className="flex items-center">
                    <div className="w-12 h-12 rounded-full bg-primary-500 text-white flex items-center justify-center mr-4">
                      {selectedMember.firstName.charAt(0)}{selectedMember.lastName.charAt(0)}
                    </div>
                    <div>
                      <p className="font-medium text-gray-900">
                        {selectedMember.firstName} {selectedMember.lastName}
                      </p>
                      <p className="text-sm text-gray-600">{selectedMember.email}</p>
                      <p className="text-sm text-gray-600">{selectedMember.phone}</p>
                      <button 
                        onClick={() => setSelectedMember(null)}
                        className="text-sm text-red-600 hover:underline mt-2"
                      >
                        Remove
                      </button>
                    </div>
                  </div>
                </div>
              )}
              
              <div className="mt-6">
                <label htmlFor="dueDate" className="label flex items-center">
                  <Calendar size={18} className="mr-2" /> Due Date
                </label>
                <input
                  id="dueDate"
                  type="date"
                  value={dueDate}
                  onChange={(e) => setDueDate(e.target.value)}
                  min={format(addDays(new Date(), 1), 'yyyy-MM-dd')}
                  className="input py-2"
                />
              </div>
            </motion.div>
          </div>
        )}
        
        {!checkoutComplete && (
          <div className="bg-gray-50 p-6 flex justify-end">
            <button
              className={`btn flex items-center ${
                selectedBook && selectedMember 
                  ? 'bg-primary-500 text-white hover:bg-primary-600' 
                  : 'bg-gray-200 text-gray-500 cursor-not-allowed'
              }`}
              disabled={!selectedBook || !selectedMember}
              onClick={handleCheckout}
            >
              Complete Checkout <ArrowRight size={18} className="ml-2" />
            </button>
          </div>
        )}
      </motion.div>
    </div>
  );
};

export default Checkout;