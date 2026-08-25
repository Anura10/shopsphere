import { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { cartAPI, orderAPI } from "../services/api";

function Checkout() {
  const navigate = useNavigate();

  const [cart, setCart] = useState(null);
  const [shippingAddress, setShippingAddress] = useState("");
  const [loading, setLoading] = useState(true);
  const [placingOrder, setPlacingOrder] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    const fetchCart = async () => {
      try {
        setLoading(true);
        setError("");

        const response = await cartAPI.get("/cart");

        const cartData =
          response.data?.cart ||
          response.data?.data ||
          response.data;

        setCart(cartData);
      } catch (err) {
        console.error(
          "Failed to load cart:",
          err
        );

        setError(
          err.response?.data?.message ||
            "Unable to load your cart."
        );
      } finally {
        setLoading(false);
      }
    };

    fetchCart();
  }, []);

  const items = cart?.items || [];

  const totalItems = items.reduce(
    (total, item) =>
      total + Number(item.quantity || 0),
    0
  );

  const cartTotal = Number(
    cart?.total || 0
  );

  const handleSubmit = async (event) => {
    event.preventDefault();

    if (placingOrder) {
      return;
    }

    const address = shippingAddress.trim();

    if (!address) {
      setError(
        "Please enter your shipping address."
      );
      return;
    }

    if (items.length === 0) {
      setError(
        "Your cart is empty. Add products before checkout."
      );
      return;
    }

    try {
      setPlacingOrder(true);
      setError("");

      const response = await orderAPI.post(
        "/orders",
        {
          shipping_address: address,
        }
      );

      console.log(
        "Order created successfully:",
        response.data
      );

      navigate("/orders");
    } catch (err) {
      console.error(
        "Failed to create order:",
        err
      );

      setError(
        err.response?.data?.message ||
          "Unable to place your order. Please try again."
      );
    } finally {
      setPlacingOrder(false);
    }
  };

  if (loading) {
    return (
      <main className="checkout-page">
        <section className="loading-state">
          <p>Loading checkout...</p>
        </section>
      </main>
    );
  }

  if (error && !cart) {
    return (
      <main className="checkout-page">
        <section className="error-state">
          <h2>Unable to load checkout</h2>

          <p>{error}</p>

          <button
            type="button"
            className="primary-button"
            onClick={() =>
              window.location.reload()
            }
          >
            Try Again
          </button>
        </section>
      </main>
    );
  }

  if (items.length === 0) {
    return (
      <main className="checkout-page">
        <section className="empty-state">
          <div className="empty-cart-icon">
            🛒
          </div>

          <h1>Your cart is empty</h1>

          <p>
            Add products to your cart before
            proceeding to checkout.
          </p>

          <Link
            to="/products"
            className="primary-button"
          >
            Browse Products
          </Link>
        </section>
      </main>
    );
  }

  return (
    <main className="checkout-page">
      <section className="checkout-header">
        <p className="section-label">
          SHOPSPHERE CHECKOUT
        </p>

        <h1>Checkout</h1>

        <p>
          Confirm your order and shipping
          information.
        </p>
      </section>

      <section className="checkout-layout">
        {/* =========================
            SHIPPING INFORMATION
        ========================== */}

        <form
          className="checkout-form"
          onSubmit={handleSubmit}
        >
          <div className="checkout-section">
            <h2>Shipping Information</h2>

            <div className="form-group">
              <label htmlFor="shippingAddress">
                Shipping Address
              </label>

              <textarea
                id="shippingAddress"
                name="shippingAddress"
                value={shippingAddress}
                onChange={(event) =>
                  setShippingAddress(
                    event.target.value
                  )
                }
                placeholder="Enter your complete shipping address"
                rows="6"
                disabled={placingOrder}
                required
              />

              <small>
                Please include your city, state,
                and PIN code.
              </small>
            </div>
          </div>

          {error && (
            <div className="error-message">
              {error}
            </div>
          )}

          <div className="checkout-actions">
            <Link
              to="/cart"
              className="secondary-button"
            >
              ← Back to Cart
            </Link>

            <button
              type="submit"
              className="primary-button"
              disabled={placingOrder}
            >
              {placingOrder
                ? "Placing Order..."
                : "Place Order"}
            </button>
          </div>
        </form>

        {/* =========================
            ORDER SUMMARY
        ========================== */}

        <aside className="checkout-summary">
          <h2>Order Summary</h2>

          <div className="checkout-items">
            {items.map((item) => {
              const name =
                item.product_name ||
                item.name ||
                "Product";

              const price = Number(
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
                  price * quantity
              );

              return (
                <div
                  className="checkout-item"
                  key={item.product_id}
                >
                  <div>
                    <strong>{name}</strong>

                    <p>
                      ₹{price.toFixed(2)} ×{" "}
                      {quantity}
                    </p>
                  </div>

                  <strong>
                    ₹{subtotal.toFixed(2)}
                  </strong>
                </div>
              );
            })}
          </div>

          <div className="summary-divider" />

          <div className="summary-row">
            <span>Total Items</span>

            <span>{totalItems}</span>
          </div>

          <div className="summary-row total-row">
            <span>Total</span>

            <strong>
              ₹{cartTotal.toFixed(2)}
            </strong>
          </div>
        </aside>
      </section>
    </main>
  );
}

export default Checkout;