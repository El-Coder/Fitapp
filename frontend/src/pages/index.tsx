import { useEffect, useState } from "react";
import { useLinkStore } from "../store/linkStore";

export default function Home() {
  const {
    fits,
    items,
    selectedFit,
    selectedItem,
    setFits,
    setItems,
    selectFit,
    selectItem,
    clearSelections,
  } = useLinkStore();

  const [linking, setLinking] = useState(false);
  const [message, setMessage] = useState("");

  // Gallery state: all linked items grouped by fit
  const [gallery, setGallery] = useState<{ fit_id: string; fit_name: string; linked_items: { item_id: string; item_name: string }[] }[]>([]);

  // Fetch gallery data from backend
  const fetchGallery = async () => {
    try {
      const res = await fetch("http://localhost:8080/api/all-linked-items");
      if (res.ok) {
        const data = await res.json();
        setGallery(data.fits || []);
      }
    } catch (e) {
      // Optionally handle error
    }
  };

  useEffect(() => {
    fetchGallery();
  }, []);

  // Fetch fits and items from backend
  useEffect(() => {
    const fetchFits = async () => {
      try {
        const res = await fetch("http://localhost:8080/api/fits");
        if (res.ok) {
          const data = await res.json();
          setFits((data.fits || []).map((fit: any) => ({ id: fit.fit_id, name: fit.fit_name })));
        }
      } catch (e) {}
    };
    const fetchItems = async () => {
      try {
        const res = await fetch("http://localhost:8080/api/items");
        if (res.ok) {
          const data = await res.json();
          setItems((data || []).map((item: any) => ({ id: item.item_id, name: item.item_name })));
        }
      } catch (e) {}
    };
    fetchFits();
    fetchItems();
  }, [setFits, setItems]);

  const handleLink = async () => {
    setLinking(true);
    setMessage("");
    try {
      // POST to your Go API: /api/link with { fit_id, item_id }
      const res = await fetch("http://localhost:8080/api/link", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ fit_id: selectedFit, item_id: selectedItem }),
      });
      if (res.ok) {
        setMessage("Item linked to fit successfully!");
        clearSelections();
        fetchGallery(); // Refresh gallery after successful link
      } else {
        const err = await res.json();
        setMessage("Error: " + (err?.error || "Unknown error"));
      }
    } catch (e) {
      setMessage("Network error");
    } finally {
      setLinking(false);
    }
  };

  return (
    <div className="p-8 max-w-md mx-auto">
      {/* Gallery of linked items */}
      <div className="mb-8">
        <h2 className="text-xl font-semibold mb-2">Linked Items Gallery</h2>
        {gallery.length === 0 ? (
          <div className="text-gray-500">No linked items yet.</div>
        ) : (
          <ul className="space-y-2">
            {gallery.map((fit) => (
  <li key={fit.fit_id}>
    <span className="font-bold">
      {fit.fit_name || fit.fit_id}:
    </span>{" "}
    {fit.linked_items && fit.linked_items.length > 0 ? (
      fit.linked_items.map(item => item.item_name || item.item_id).join(", ")
    ) : (
      <span className="text-gray-400">No linked items</span>
    )}
  </li>
))}
          </ul>
        )}
      </div>
      <h1 className="text-2xl font-bold mb-4">Link Item to Fit</h1>
      <div className="mb-4">
        <label className="block mb-1 font-semibold" htmlFor="fit-select">
          Select Fit
        </label>
        <select
          id="fit-select"
          className="w-full p-2 border rounded mb-2"
          value={selectedFit}
          onChange={(e) => selectFit(e.target.value)}
        >
          <option value="">-- Choose a fit --</option>
          {gallery.map((fit) => (
            <option key={fit.fit_id} value={fit.fit_id}>
              {fit.fit_name || fit.fit_id}
            </option>
          ))}
        </select>
      </div>
      <div className="mb-4">
        <label className="block mb-1 font-semibold" htmlFor="item-select">
          Select Item
        </label>
        <select
          id="item-select"
          className="w-full p-2 border rounded mb-2"
          value={selectedItem}
          onChange={(e) => selectItem(e.target.value)}
        >
          <option value="">-- Choose an item --</option>
          {items.map((item) => (
            <option key={item.id} value={item.id}>
              {item.name}
            </option>
          ))}
        </select>
      </div>
      <button
        className="bg-blue-600 text-white px-4 py-2 rounded disabled:opacity-50 w-full"
        onClick={handleLink}
        disabled={linking || !selectedFit || !selectedItem}
      >
        {linking ? "Linking..." : "Link"}
      </button>
      {message && <div className="mt-4 text-green-600">{message}</div>}
    </div>
  );
}