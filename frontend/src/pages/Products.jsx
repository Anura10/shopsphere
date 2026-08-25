import { useEffect, useMemo, useState } from "react";
import { Link } from "react-router-dom";
import { productAPI } from "../services/api";

function Products() {
  const [products, setProducts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const [searchTerm, setSearchTerm] = useState("");
  const [categoryFilter, setCategoryFilter] = useState("all");
  const [sortOption, setSortOption] = useState("default");

  const fetchProducts = async () => {
    try {
      setLoading(true);
      setError("");

      const response =
        await productAPI.get("/products");

      setProducts(
        response.data?.products || []
      );
    } catch (err) {
      console.error(
        "Failed to fetch products:",
        err
      );

      setError(
        err.response?.data?.message ||
          "Unable to load products."
      );
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchProducts();
  }, []);

  const categories = useMemo(() => {
    const uniqueCategories = [
      ...new Set(
        products
          .map((product) => product.category_id)
          .filter(
            (category) =>
              category !== null &&
              category !== undefined
          )
      ),
    ];

    return uniqueCategories.sort(
      (a, b) => Number(a) - Number(b)
    );
  }, [products]);

  const filteredProducts = useMemo(() => {
    const search = searchTerm
      .trim()
      .toLowerCase();

    let result = products.filter((product) => {
      const matchesSearch =
        !search ||
        product.name
          ?.toLowerCase()
          .includes(search) ||
        product.description
          ?.toLowerCase()
          .includes(search) ||
        product.sku
          ?.toLowerCase()
          .includes(search);

      const matchesCategory =
        categoryFilter === "all" ||
        String(product.category_id) ===
          String(categoryFilter);

      return (
        matchesSearch &&
        matchesCategory
      );
    });

    result = [...result].sort(
      (a, b) => {
        switch (sortOption) {
          case "price-low":
            return (
              Number(a.price) -
              Number(b.price)
            );

          case "price-high":
            return (
              Number(b.price) -
              Number(a.price)
            );

          case "name-asc":
            return (a.name || "").localeCompare(
              b.name || ""
            );

          case "name-desc":
            return (b.name || "").localeCompare(
              a.name || ""
            );

          default:
            return 0;
        }
      }
    );

    return result;
  }, [
    products,
    searchTerm,
    categoryFilter,
    sortOption,
  ]);

  const clearFilters = () => {
    setSearchTerm("");
    setCategoryFilter("all");
    setSortOption("default");
  };

  if (loading) {
    return (
      <main className="products-page">
        <section className="products-header">
          <p className="section-label">
            SHOPSPHERE STORE
          </p>

          <h1>All Products</h1>
        </section>

        <div className="loading-state">
          Loading products...
        </div>
      </main>
    );
  }

  if (error) {
    return (
      <main className="products-page">
        <section className="products-header">
          <p className="section-label">
            SHOPSPHERE STORE
          </p>

          <h1>All Products</h1>
        </section>

        <div className="error-message">
          {error}
        </div>

        <button
          type="button"
          className="primary-button"
          onClick={fetchProducts}
        >
          Try Again
        </button>
      </main>
    );
  }

  return (
    <main className="products-page">
      <section className="products-header">
        <p className="section-label">
          SHOPSPHERE STORE
        </p>

        <h1>All Products</h1>

        <p>
          Discover products available in the
          ShopSphere catalog.
        </p>
      </section>

      {/* =========================
          PRODUCT CONTROLS
      ========================== */}

      <section className="product-controls">
        <div className="search-control">
          <label htmlFor="product-search">
            Search Products
          </label>

          <input
            id="product-search"
            type="search"
            value={searchTerm}
            onChange={(event) =>
              setSearchTerm(
                event.target.value
              )
            }
            placeholder="Search by name, description or SKU..."
          />
        </div>

        <div className="filter-control">
          <label htmlFor="category-filter">
            Category
          </label>

          <select
            id="category-filter"
            value={categoryFilter}
            onChange={(event) =>
              setCategoryFilter(
                event.target.value
              )
            }
          >
            <option value="all">
              All Categories
            </option>

            {categories.map((category) => (
              <option
                key={category}
                value={category}
              >
                Category #{category}
              </option>
            ))}
          </select>
        </div>

        <div className="filter-control">
          <label htmlFor="sort-products">
            Sort By
          </label>

          <select
            id="sort-products"
            value={sortOption}
            onChange={(event) =>
              setSortOption(
                event.target.value
              )
            }
          >
            <option value="default">
              Default
            </option>

            <option value="price-low">
              Price: Low to High
            </option>

            <option value="price-high">
              Price: High to Low
            </option>

            <option value="name-asc">
              Name: A to Z
            </option>

            <option value="name-desc">
              Name: Z to A
            </option>
          </select>
        </div>

        <button
          type="button"
          className="secondary-button"
          onClick={clearFilters}
        >
          Clear Filters
        </button>
      </section>

      {/* =========================
          RESULT COUNT
      ========================== */}

      <section className="products-result-info">
        <p>
          Showing{" "}
          <strong>
            {filteredProducts.length}
          </strong>{" "}
          of{" "}
          <strong>
            {products.length}
          </strong>{" "}
          product
          {products.length !== 1
            ? "s"
            : ""}
        </p>
      </section>

      {/* =========================
          PRODUCTS
      ========================== */}

      {filteredProducts.length === 0 ? (
        <section className="empty-state">
          <h2>
            No products found
          </h2>

          <p>
            Try changing your search or
            filters.
          </p>

          <button
            type="button"
            className="primary-button"
            onClick={clearFilters}
          >
            Clear Filters
          </button>
        </section>
      ) : (
        <section className="product-grid">
          {filteredProducts.map(
            (product) => {
              const stock = Number(
                product.stock_quantity || 0
              );

              const isOutOfStock =
                stock <= 0;

              const isLowStock =
                stock > 0 && stock <= 5;

              return (
                <article
                  className="product-card"
                  key={product.id}
                >
                  <div className="product-image">
                    {product.image_url ? (
                      <img
                        src={product.image_url}
                        alt={product.name}
                        loading="lazy"
                      />
                    ) : (
                      <span>
                        ShopSphere
                      </span>
                    )}
                  </div>

                  <div className="product-content">
                    <p className="product-category">
                      Category #
                      {product.category_id}
                    </p>

                    <h2>
                      {product.name}
                    </h2>

                    <p className="product-description">
                      {product.description ||
                        "Quality product from ShopSphere."}
                    </p>

                    <p className="product-sku">
                      SKU: {product.sku}
                    </p>

                    <div className="product-footer">
                      <strong>
                        ₹
                        {Number(
                          product.price
                        ).toFixed(2)}
                      </strong>

                      <span
                        className={
                          isOutOfStock
                            ? "stock-out"
                            : isLowStock
                            ? "stock-low"
                            : "stock-available"
                        }
                      >
                        {isOutOfStock
                          ? "Out of stock"
                          : isLowStock
                          ? `Only ${stock} left`
                          : `In stock: ${stock}`}
                      </span>
                    </div>

                    <Link
                      to={`/products/${product.id}`}
                      className="primary-button"
                    >
                      View Product
                    </Link>
                  </div>
                </article>
              );
            }
          )}
        </section>
      )}
    </main>
  );
}

export default Products;