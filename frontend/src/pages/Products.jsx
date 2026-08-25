import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { productAPI } from "../services/api";

function Products() {
  const [products, setProducts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const fetchProducts = async () => {
    try {
      setLoading(true);
      setError("");

      const response =
        await productAPI.get("/products");

      setProducts(
        response.data.products || []
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

        <p>
          {products.length} product
          {products.length !== 1
            ? "s"
            : ""}{" "}
          available
        </p>
      </section>

      {products.length === 0 ? (
        <section className="empty-state">
          <h2>
            No products found
          </h2>

          <p>
            There are currently no
            products available.
          </p>
        </section>
      ) : (
        <section className="product-grid">
          {products.map((product) => (
            <article
              className="product-card"
              key={product.id}
            >
              <div className="product-image">
                {product.image_url ? (
                  <img
                    src={product.image_url}
                    alt={product.name}
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

                  <span>
                    Stock:{" "}
                    {product.stock_quantity}
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
          ))}
        </section>
      )}
    </main>
  );
}

export default Products;