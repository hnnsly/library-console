import { env } from '$env/dynamic/private';

// Определим типы для API ответов (хорошей практикой будет вынести их в отдельный src/lib/types.ts)
export interface Pharmacy {
	id: number;
	name: string;
	created_at: string;
	updated_at: string;
}

export interface Product {
	id: number;
	name: string;
	dosages: string[];
	created_at: string;
	updated_at: string;
}

// ... другие типы для InventoryItem, Supplier, Order и т.д.

async function request<T>(
	method: 'GET' | 'POST' | 'PUT' | 'DELETE',
	endpoint: string,
	data?: unknown,
	customHeaders?: HeadersInit
): Promise<T> {
	const headers: HeadersInit = {
		'Content-Type': 'application/json',
		...customHeaders
	};

	// Добавляем токен для POST, PUT, DELETE запросов
	// ВАЖНО: Убедитесь, что ваш API действительно требует токен для всех этих методов.
	// Если какой-то GET-запрос тоже требует токен, логику нужно будет скорректировать.
	if (['POST', 'PUT', 'DELETE'].includes(method)) {
		if (!env.AUTH_TOKEN) {
			console.warn('AUTH_TOKEN is not set for a protected API call.');
			// В реальном приложении здесь может быть более строгая обработка ошибки
		}
		const authHeaders: HeadersInit = {
			...headers,
			Authorization: `Bearer ${env.AUTH_TOKEN}`
		};
		return request<T>(method, endpoint, data, authHeaders);
	}

	const url = `${env.API_BASE_URL}${endpoint}`;
	console.log(`API Request: ${method} ${url}`); // Для отладки

	try {
		const response = await fetch(url, {
			method,
			headers,
			body: data ? JSON.stringify(data) : undefined
		});

		if (!response.ok) {
			const errorBody = await response.text();
			console.error(`API Error: ${response.status} ${response.statusText}`, errorBody);
			throw new Error(
				`API request failed: ${response.status} ${response.statusText} - ${errorBody}`
			);
		}

		if (response.status === 204) {
			// No Content
			return undefined as T;
		}

		return response.json() as T;
	} catch (error) {
		console.error('Fetch API error:', error);
		throw error; // Перебрасываем ошибку для обработки выше
	}
}

// --- Pharmacies ---
export const getPharmacies = (limit = 10, offset = 0): Promise<Pharmacy[]> => {
	return request<Pharmacy[]>('GET', `/pharmacies/?limit=${limit}&offset=${offset}`);
};

export const getPharmacyById = (id: number | string): Promise<Pharmacy> => {
	return request<Pharmacy>('GET', `/pharmacies/${id}`);
};

export const createPharmacy = (name: string): Promise<Pharmacy> => {
	return request<Pharmacy>('POST', '/pharmacies/', { name });
};

export const updatePharmacy = (id: number | string, name: string): Promise<Pharmacy> => {
	return request<Pharmacy>('PUT', `/pharmacies/${id}`, { name });
};

export const deletePharmacy = (id: number | string): Promise<void> => {
	return request<void>('DELETE', `/pharmacies/${id}`);
};

// --- Products ---
export const getProducts = (limit = 10, offset = 0): Promise<Product[]> => {
	return request<Product[]>('GET', `/products/?limit=${limit}&offset=${offset}`);
};

// ... Добавьте здесь другие функции для работы с Products, Inventory, Suppliers, Orders, Stock Search
// Например:
// export const addInventoryItem = (pharmacyId: string, itemData: any) => { ... }
// export const createOrder = (orderData: any) => { ... }
