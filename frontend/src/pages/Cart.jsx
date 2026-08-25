import { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { cartAPI } from "../services/api";

function Cart() {
  const navigate = useNavigate();

  const [cart, setCart] = useState(null);
  const [loading, setLoading] = useState(true);
  const [actionLoading, setActionLoading] = useState(false);
  const [error, setError] = useState("");

  const fetchCart = async () => {
    try {
      setLoading(true);
      setError("");

      const response = await cartAPI.get("/cart");

      setCart(
        response.data?.cart ||
          response.data?.data ||
          response.data
      );
    } catch (err) {
      console.error("Failed to fetch cart:", err);

      setError(
        err.response?.data?.message ||
          "Unable to load your cart."
      );
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchCart();
  }, []);

  const updateQuantity = async (
    productId,
    quantity
  ) => {
    if (quantity < 1 || actionLoading) {
      return;
    }

    try {
      setActionLoading(true);
      setError("");

      await cartAPI.put(
        `/cart/items/${productId}`,
        {
          quantity,
        }
      );

      await fetchCart();
    } catch (err) {
      console.error(
        "Failed to update quantity:",
        err
      );

      setError(
        err.response?.data?.message ||
          "Unable to update cart."
      );
    } finally {
      setActionLoading(false);
    }
  };

  const removeItem = async (productId) => {
    if (actionLoading) {
      return;
    }

    try {
      setActionLoading(true);
      setError("");

      await cartAPI.delete(
        `/cart/items/${productId}`
      );

      await fetchCart();
    } catch (err) {
      console.error(
        "Failed to remove item:",
        err
      );

      setError(
        err.response?.data?.message ||
          "Unable to remove item."
      );
    } finally {
      setActionLoading(false);
    }
  };

  const clearCart = async () => {
    if (actionLoading) {
      return;
    }

    const confirmed = window.confirm(
      "Are you sure you want to clear your cart?"
    );

    if (!confirmed) {
      return;
    }

    try {
      setActionLoading(true);
      setError("");

      await cartAPI.delete("/cart");

      await fetchCart();
    } catch (err) {
      console.error(
        "Failed to clear cart:",
        err
      );

      setError(
        err.response?.data?.message ||
          "Unable to clear cart."
      );
    } finally {
      setActionLoading(false);
    }
  };

  if (loading) {
    return (
      <main className="cart-page">
        <section className="loading-state">
          <p>Loading your cart...</p>
        </section>
      </main>
    );
  }

  if (error && !cart) {
    return (
      <main className="cart-page">
        <section className="error-state">
          <h2>Unable to load cart</h2>

          <p>{error}</p>

          <button
            type="button"
            className="primary-button"
            onClick={fetchCart}
          >
            Try Again
          </button>
        </section>
      </main>
    );
  }

  const items = cart?.items || [];

  const totalItems = items.reduce(
    (total, item) =>
      total + Number(item.quantity || 0),
    0
  );

  const cartTotal = Number(
    cart?.total || 0
  );

  return (
    <main className="cart-page">
      <section className="cart-header">
        <p className="section-label">
          SHOPSPHERE CART
        </p>

        <h1>Your Shopping Cart</h1>

        <p>
          Review your items before checkout.
        </p>
      </section>

      {error && (
        <div className="error-message">
          {error}
        </div>
      )}

      {items.length === 0 ? (
        <section className="empty-state">
          <div className="empty-cart-icon">
            🛒
          </div>

          <h2>Your cart is empty</h2>

          <p>
            You haven't added any products yet.
          </p>

          <Link
            to="/products"
            className="primary-button"
          >
            Continue Shopping
          </Link>
        </section>
      ) : (
        <section className="cart-layout">
          <div className="cart-items">
            {items.map((item) => {
              const productName =
                item.product_name ||
                item.name ||
                "Product";

              const productPrice = Number(
                item.product_price ||
                  item.price ||
                  item.unit_price ||
                  0
              );

              const quantity = Number(
                item.quantity || 0
              );

              const subtotal = Number(
                item.subtotal ||
                  productPrice * quantity
              );

              return (
                <article
                  className="cart-item"
                  key={item.product_id}
                >
                  <div className="cart-item-info">
                    <h2>{productName}</h2>

                    <p>
                      ₹
                      {productPrice.toFixed(2)}
                      {" "}each
                    </p>
                  </div>

                  <div className="cart-item-actions">
                    <div className="quantity-control">
                      <button
                        type="button"
                        aria-label={`Decrease quantity of ${productName}`}
                        onClick={() =>
                          updateQuantity(
                            item.product_id,
                            quantity - 1
                          )
                        }
                        disabled={
                          quantity <= 1 ||
                          actionLoading
                        }
                      >
                        −
                      </button>

                      <span>
                        {quantity}
                      </span>

                      <button
                        type="button"
                        aria-label={`Increase quantity of ${productName}`}
                        onClick={() =>
                          updateQuantity(
                            item.product_id,
                            quantity + 1
                          )
                        }
                        disabled={actionLoading}
                      >
                        +
                      </button>
                    </div>

                    <button
                      type="button"
                      className="remove-button"
                      onClick={() =>
                        removeItem(
                          item.product_id
                        )
                      }
                      disabled={actionLoading}
                    >
                      Remove
                    </button>
                  </div>

                  <strong className="item-subtotal">
                    ₹{subtotal.toFixed(2)}
                  </strong>
                </article>
              );
            })}
          </div>

          <aside className="cart-summary">
            <h2>Order Summary</h2>

            <div className="summary-row">
              <span>Items</span>

              <span>{totalItems}</span>
            </div>

            <div className="summary-row">
              <span>Subtotal</span>

              <span>
                ₹{cartTotal.toFixed(2)}
              </span>
            </div>

            <div className="summary-divider" />

            <div className="summary-row total-row">
              <span>Total</span>

              <strong>
                ₹{cartTotal.toFixed(2)}
              </strong>
            </div>

            <button
              type="button"
              className="primary-button"
              onClick={() =>
                navigate("/checkout")
              }
              disabled={actionLoading}
            >
              Proceed to Checkout
            </button>

            <button
              type="button"
              className="secondary-button"
              onClick={clearCart}
              disabled={actionLoading}
            >
              Clear Cart
            </button>

            <Link
              to="/products"
              className="continue-shopping"
            >
              ← Continue Shopping
            </Link>
          </aside>
        </section>
      )}
    </main>
  );
}

export default Cart;