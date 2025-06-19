import { create } from "zustand";

type Fit = { id: string; name: string };
type Item = { id: string; name: string };

type LinkStore = {
  fits: Fit[];
  items: Item[];
  selectedFit: string;
  selectedItem: string;
  setFits: (fits: Fit[]) => void;
  setItems: (items: Item[]) => void;
  selectFit: (id: string) => void;
  selectItem: (id: string) => void;
  clearSelections: () => void;
};

export const useLinkStore = create<LinkStore>((set) => ({
  fits: [],
  items: [],
  selectedFit: "",
  selectedItem: "",
  setFits: (fits) => set({ fits }),
  setItems: (items) => set({ items }),
  selectFit: (id) => set({ selectedFit: id }),
  selectItem: (id) => set({ selectedItem: id }),
  clearSelections: () => set({ selectedFit: "", selectedItem: "" }),
}));