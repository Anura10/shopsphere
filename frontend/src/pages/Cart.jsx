import { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { cartAPI } from "../services/api";

function Cart() {
  const navigate = useNavigate();

  const [cart, setCart] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const fetchCart = async () => {
    try {
      setLoading(true);
      setError("");

      const response = await cartAPI.get("/cart");

      setCart(
        response.data.cart ||
          response.data.data ||
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
    if (quantity <= 0) {
      return;
    }

    try {
      await cartAPI.put(
        `/cart/items/${productId}`,
        {
          quantity,
        }
      );

      await fetchCart();
    } catch (err) {
      setError(
        err.response?.data?.message ||
          "Unable to update cart."
      );
    }
  };

  const removeItem = async (productId) => {
    try {
      await cartAPI.delete(
        `/cart/items/${productId}`
      );

      await fetchCart();
    } catch (err) {
      setError(
        err.response?.data?.message ||
          "Unable to remove item."
      );
    }
  };

  const clearCart = async () => {
    try {
      await cartAPI.delete("/cart");

      await fetchCart();
    } catch (err) {
      setError(
        err.response?.data?.message ||
          "Unable to clear cart."
      );
    }
  };

  if (loading) {
    return (
      <main className="cart-page">
        <div className="loading-state">
          Loading cart...
        </div>
      </main>
    );
  }

  if (error && !cart) {
    return (
      <main className="cart-page">
        <div className="error-message">
          {error}
        </div>

        <button
          type="button"
          onClick={fetchCart}
        >
          Try Again
        </button>
      </main>
    );
  }

  const items = cart?.items || [];

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
          <h2>Your cart is empty</h2>

          <p>
            Add some products to your cart
            before checking out.
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
            {items.map((item) => (
              <article
                className="cart-item"
                key={item.product_id}
              >
                <div className="cart-item-info">
                  <h2>
                    {item.product_name ||
                      item.name}
                  </h2>

                  <p>
                    ₹
                    {Number(
                      item.product_price ||
                        item.price ||
                        item.unit_price ||
                        0
                    ).toFixed(2)}
                    {" "}each
                  </p>
                </div>

                <div className="cart-item-actions">
                  <button
                    type="button"
                    onClick={() =>
                      updateQuantity(
                        item.product_id,
                        item.quantity - 1
                      )
                    }
                    disabled={item.quantity <= 1}
                  >
                    −
                  </button>

                  <span>
                    {item.quantity}
                  </span>

                  <button
                    type="button"
                    onClick={() =>
                      updateQuantity(
                        item.product_id,
                        item.quantity + 1
                      )
                    }
                  >
                    +
                  </button>

                  <button
                    type="button"
                    onClick={() =>
                      removeItem(
                        item.product_id
                      )
                    }
                  >
                    Remove
                  </button>
                </div>

                <strong>
                  ₹
                  {Number(
                    item.subtotal || 0
                  ).toFixed(2)}
                </strong>
              </article>
            ))}
          </div>

          <aside className="cart-summary">
            <h2>Order Summary</h2>

            <div className="summary-row">
              <span>Items</span>

              <span>
                {items.reduce(
                  (total, item) =>
                    total + item.quantity,
                  0
                )}
              </span>
            </div>

            <div className="summary-row total-row">
              <span>Total</span>

              <strong>
                ₹
                {Number(
                  cart?.total || 0
                ).toFixed(2)}
              </strong>
            </div>

            <button
              type="button"
              className="primary-button"
              onClick={() =>
                navigate("/checkout")
              }
            >
              Proceed to Checkout
            </button>

            <button
              type="button"
              className="secondary-button"
              onClick={clearCart}
            >
              Clear Cart
            </button>

            <Link
              to="/products"
              className="continue-shopping"
            >
              Continue Shopping
            </Link>
          </aside>
        </section>
      )}
    </main>
  );
}

export default Cart;