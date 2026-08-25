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
      setMessage("");

      const response = await productAPI.get(
        `/products/${productId}`
      );

      const data = response.data;

      setProduct(
        data?.product ||
          data?.data ||
          data
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

  const stock = Number(
    product?.stock_quantity || 0
  );

  const price = Number(
    product?.price || 0
  );

  const isOutOfStock = stock <= 0;
  const isLowStock =
    stock > 0 && stock <= 5;

  const increaseQuantity = () => {
    if (quantity < stock) {
      setQuantity(
        (current) => current + 1
      );
    }
  };

  const decreaseQuantity = () => {
    if (quantity > 1) {
      setQuantity(
        (current) => current - 1
      );
    }
  };

  const handleQuantityChange = (
    event
  ) => {
    const value = Number(
      event.target.value
    );

    if (!Number.isInteger(value)) {
      return;
    }

    if (value < 1) {
      setQuantity(1);
      return;
    }

    if (value > stock) {
      setQuantity(stock);
      return;
    }

    setQuantity(value);
  };

  const handleAddToCart = async () => {
    if (isOutOfStock) {
      setError(
        "This product is currently out of stock."
      );
      return;
    }

    if (quantity < 1 || quantity > stock) {
      setError(
        "Please select a valid quantity."
      );
      return;
    }

    try {
      setAdding(true);
      setError("");
      setMessage("");

      await cartAPI.post(
        "/cart/items",
        {
          product_id: Number(product.id),
          quantity,
        }
      );

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

  const handleBuyNow = async () => {
    if (isOutOfStock) {
      setError(
        "This product is currently out of stock."
      );
      return;
    }

    if (quantity < 1 || quantity > stock) {
      setError(
        "Please select a valid quantity."
      );
      return;
    }

    try {
      setAdding(true);
      setError("");
      setMessage("");

      await cartAPI.post(
        "/cart/items",
        {
          product_id: Number(product.id),
          quantity,
        }
      );

      navigate("/checkout");
    } catch (err) {
      console.error(
        "Failed to buy product:",
        err
      );

      setError(
        err.response?.data?.message ||
          "Unable to continue to checkout."
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

        <Link
          to="/products"
          className="secondary-button"
        >
          ← Back to Products
        </Link>
      </main>
    );
  }

  if (!product) {
    return (
      <main className="product-details-page">
        <section className="empty-state">
          <h1>Product not found</h1>

          <p>
            The product you're looking for
            doesn't exist.
          </p>

          <Link
            to="/products"
            className="primary-button"
          >
            ← Back to Products
          </Link>
        </section>
      </main>
    );
  }

  const totalPrice =
    price * quantity;

  return (
    <main className="product-details-page">
      <Link
        to="/products"
        className="back-link"
      >
        ← Back to Products
      </Link>

      <section className="product-details-layout">
        {/* =========================
            PRODUCT IMAGE
        ========================== */}

        <div className="product-details-image">
          {product.image_url ? (
            <img
              src={product.image_url}
              alt={product.name}
            />
          ) : (
            <div className="product-image-placeholder">
              ShopSphere
            </div>
          )}
        </div>

        {/* =========================
            PRODUCT INFORMATION
        ========================== */}

        <div className="product-details-content">
          <p className="product-category">
            Category #{product.category_id}
          </p>

          <h1>{product.name}</h1>

          <p className="product-price">
            ₹{price.toFixed(2)}
          </p>

          <p className="product-description">
            {product.description ||
              "Quality product from ShopSphere."}
          </p>

          <p className="product-sku">
            SKU: {product.sku}
          </p>

          {/* =========================
              STOCK STATUS
          ========================== */}

          <div className="product-stock">
            {isOutOfStock ? (
              <span className="stock-out">
                Out of stock
              </span>
            ) : isLowStock ? (
              <span className="stock-low">
                Only {stock} left in stock
              </span>
            ) : (
              <span className="stock-available">
                In stock: {stock}
              </span>
            )}
          </div>

          {/* =========================
              QUANTITY
          ========================== */}

          {!isOutOfStock && (
            <div className="quantity-section">
              <label htmlFor="product-quantity">
                Quantity
              </label>

              <div className="quantity-control">
                <button
                  type="button"
                  onClick={decreaseQuantity}
                  disabled={
                    quantity <= 1 ||
                    adding
                  }
                  aria-label="Decrease quantity"
                >
                  −
                </button>

                <input
                  id="product-quantity"
                  type="number"
                  min="1"
                  max={stock}
                  value={quantity}
                  onChange={
                    handleQuantityChange
                  }
                  disabled={adding}
                />

                <button
                  type="button"
                  onClick={increaseQuantity}
                  disabled={
                    quantity >= stock ||
                    adding
                  }
                  aria-label="Increase quantity"
                >
                  +
                </button>
              </div>
            </div>
          )}

          {/* =========================
              TOTAL
          ========================== */}

          {!isOutOfStock && (
            <div className="product-total">
              <span>Total</span>

              <strong>
                ₹{totalPrice.toFixed(2)}
              </strong>
            </div>
          )}

          {/* =========================
              MESSAGES
          ========================== */}

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

          {/* =========================
              ACTIONS
          ========================== */}

          <div className="product-actions">
            <button
              type="button"
              className="primary-button"
              onClick={handleAddToCart}
              disabled={
                adding || isOutOfStock
              }
            >
              {adding
                ? "Adding..."
                : isOutOfStock
                ? "Out of Stock"
                : "Add to Cart"}
            </button>

            <button
              type="button"
              className="secondary-button"
              onClick={handleBuyNow}
              disabled={
                adding || isOutOfStock
              }
            >
              {adding
                ? "Processing..."
                : "Buy Now"}
            </button>
          </div>

          <Link
            to="/cart"
            className="continue-shopping"
          >
            View Cart
          </Link>
        </div>
      </section>
    </main>
  );
}

export default ProductDetails;