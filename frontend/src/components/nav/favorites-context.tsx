import { type ReactNode, createContext, useContext, useState } from "react";
import { toast } from "sonner";
import type { NavItemDataProps } from "./types";

interface FavoritesContextType {
	favorites: NavItemDataProps[];
	addFavorite: (item: NavItemDataProps) => void;
	removeFavorite: (path: string) => void;
	isFavorite: (path: string) => boolean;
}

const FavoritesContext = createContext<FavoritesContextType | undefined>(undefined);

export function FavoritesProvider({ children }: { children: ReactNode }) {
	const [favorites, setFavorites] = useState<NavItemDataProps[]>(() => {
		const saved = localStorage.getItem("navigation-favorites");
		return saved ? JSON.parse(saved) : [];
	});

	const addFavorite = (item: NavItemDataProps) => {
		if (favorites.length >= 8) {
			toast.error("Add up to 8 favorites.");
			return;
		}
		setFavorites((prev) => {
			const newFavorites = [...prev, item];
			localStorage.setItem("navigation-favorites", JSON.stringify(newFavorites));
			return newFavorites;
		});
	};

	const removeFavorite = (path: string) => {
		setFavorites((prev) => {
			const newFavorites = prev.filter((item) => item.path !== path);
			localStorage.setItem("navigation-favorites", JSON.stringify(newFavorites));
			return newFavorites;
		});
	};

	const isFavorite = (path: string) => {
		return favorites.some((item) => item.path === path);
	};

	return (
		<FavoritesContext.Provider value={{ favorites, addFavorite, removeFavorite, isFavorite }}>
			{children}
		</FavoritesContext.Provider>
	);
}

export function useFavorites() {
	const context = useContext(FavoritesContext);
	if (context === undefined) {
		throw new Error("useFavorites must be used within a FavoritesProvider");
	}
	return context;
}
