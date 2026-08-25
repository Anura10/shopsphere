import { useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import { productAPI, cartAPI } from "../services/api";

function ProductDetails() {
  const { productId } = useParams();
  const navigate = useNavigate();

  const [product, setProduct] = useState(null);
  const [quantity, setQuantity] = useState(1);

  const [loading, setLoading] = useState(true);
  const [adding, setAdding] = useState(false);

  const [error, setError] = useState("");
  const [message, setMessage] = useState("");

  const fetchProduct = async () => {
    try {
      setLoading(true);
      setError("");

      const response = await productAPI.get(
        `/products/${productId}`
      );

      const data = response.data;

      setProduct(
        data.product || data.data || data
      );
    } catch (err) {
      console.error(
        "Failed to fetch product:",
        err
      );

      setError(
        err.response?.data?.message ||
          "Unable to load product."
      );
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchProduct();
  }, [productId]);

  const increaseQuantity = () => {
    if (
      product &&
      quantity < product.stock_quantity
    ) {
      setQuantity((current) => current + 1);
    }
  };

  const decreaseQuantity = () => {
    if (quantity > 1) {
      setQuantity((current) => current - 1);
    }
  };

  const handleAddToCart = async () => {
    try {
      setAdding(true);
      setError("");
      setMessage("");

      await cartAPI.post("/cart/items", {
        product_id: Number(product.id),
        quantity,
      });

      setMessage(
        "Product added to your cart successfully."
      );
    } catch (err) {
      console.error(
        "Failed to add product to cart:",
        err
      );

      setError(
        err.response?.data?.message ||
          "Unable to add product to cart."
      );
    } finally {
      setAdding(false);
    }
  };

  if (loading) {
    return (
      <main className="product-details-page">
        <div className="loading-state">
          Loading product...
        </div>
      </main>
    );
  }

  if (error && !product) {
    return (
      <main className="product-details-page">
        <div className="error-message">
          {error}
        </div>

        <Link to="/products">
          Back to Products
        </Link>
      </main>
    );
  }

  if (!product) {
    return (
      <main className="product-details-page">
        <h1>Product not found</h1>

        <Link to="/products">
          Back to Products
        </Link>
      </main>
    );
  }

  const totalPrice =
    Number(product.price) * quantity;

  return (
    <main className="product-details-page">
      <Link
        to="/products"
        className="back-link"
      >
        ← Back to Products
      </Link>

      <section className="product-details">
        <div className="product-details-image">
          {product.image_url ? (
            <img
              src={product.image_url}
              alt={product.name}
            />
          ) : (
            <span>ShopSphere</span>
          )}
        </div>

        <div className="product-details-content">
          <p className="section-label">
            SHOPSPHERE PRODUCT
          </p>

          <h1>{product.name}</h1>

          <p className="product-sku">
            SKU: {product.sku}
          </p>

          <p className="product-details-description">
            {product.description}
          </p>

          <div className="product-price">
            ₹{Number(product.price).toFixed(2)}
          </div>

          <div className="stock-info">
            {product.stock_quantity > 0 ? (
              <span>
                {product.stock_quantity} available
              </span>
            ) : (
              <span>Out of stock</span>
            )}
          </div>

          {error && (
            <div className="error-message">
              {error}
            </div>
          )}

          {message && (
            <div className="success-message">
              {message}
            </div>
          )}

          {product.stock_quantity > 0 && (
            <>
              <div className="quantity-selector">
                <button
                  type="button"
                  onClick={decreaseQuantity}
                  disabled={quantity <= 1}
                >
                  −
                </button>

                <span>{quantity}</span>

                <button
                  type="button"
                  onClick={increaseQuantity}
                  disabled={
                    quantity >=
                    product.stock_quantity
                  }
                >
                  +
                </button>
              </div>

              <p>
                Total:{" "}
                <strong>
                  ₹{totalPrice.toFixed(2)}
                </strong>
              </p>

              <button
                type="button"
                className="primary-button"
                onClick={handleAddToCart}
                disabled={adding}
              >
                {adding
                  ? "Adding..."
                  : "Add to Cart"}
              </button>

              {message && (
                <button
                  type="button"
                  className="secondary-button"
                  onClick={() => navigate("/cart")}
                >
                  View Cart
                </button>
              )}
            </>
          )}
        </div>
      </section>
    </main>
  );
}

export default ProductDetails;