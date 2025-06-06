import React from "react";
import { useNavigate } from "react-router-dom";
import { BookX } from "lucide-react";

const NotFound: React.FC = () => {
  const navigate = useNavigate();

  return (
    <div className="min-h-[80vh] flex flex-col items-center justify-center text-center p-4">
      <BookX size={64} className="text-primary-500 mb-6" />
      <h1 className="text-4xl font-serif font-bold text-gray-900 mb-4">
        Страница не найдена
      </h1>
      <p className="text-xl text-gray-600 mb-8 max-w-md">
        Страница, которую вы ищете, не найдена.
      </p>
      <button onClick={() => navigate("/")} className="btn btn-primary">
        Вернуться на Главную
      </button>
    </div>
  );
};

export default NotFound;
