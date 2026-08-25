import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { orderAPI } from "../services/api";

function Orders() {
  const navigate = useNavigate();

  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    const fetchOrders = async () => {
      try {
        setLoading(true);
        setError("");

        const response = await orderAPI.get("/orders");

        setOrders(response.data?.orders || []);
      } catch (err) {
        console.error("Failed to fetch orders:", err);

        setError(
          err.response?.data?.message ||
            "Unable to load your orders."
        );
      } finally {
        setLoading(false);
      }
    };

    fetchOrders();
  }, []);

  if (loading) {
    return (
      <main className="orders-page">
        <h1>My Orders</h1>
        <p>Loading orders...</p>
      </main>
    );
  }

  return (
    <main className="orders-page">
      <section className="orders-header">
        <p className="section-label">
          SHOPSPHERE
        </p>

        <h1>My Orders</h1>

        <p>
          View your previous ShopSphere orders.
        </p>
      </section>

      {error && (
        <div className="error-message">
          {error}
        </div>
      )}

      {!error && orders.length === 0 && (
        <section className="empty-orders">
          <h2>No orders yet</h2>

          <p>
            Your completed orders will appear here.
          </p>

          <button
            className="primary-button"
            onClick={() => navigate("/products")}
          >
            Start Shopping
          </button>
        </section>
      )}

      <section className="orders-list">
        {orders.map((order) => (
          <article
            className="order-card"
            key={order.id}
          >
            <div className="order-card-header">
              <div>
                <h2>Order #{order.id}</h2>

                <p>
                  {order.created_at
                    ? new Date(
                        order.created_at
                      ).toLocaleString()
                    : "Date unavailable"}
                </p>
              </div>

              <span className="order-status">
                {order.status}
              </span>
            </div>

            <div className="order-card-body">
              <p>
                <strong>Total:</strong>{" "}
                ₹{Number(order.total_amount).toFixed(2)}
              </p>

              <p>
                <strong>Shipping:</strong>{" "}
                {order.shipping_address}
              </p>
            </div>

            {order.items?.length > 0 && (
              <div className="order-items">
                <h3>Items</h3>

                {order.items.map((item) => (
                  <div
                    className="order-item"
                    key={item.id}
                  >
                    <span>
                      {item.product_name}
                    </span>

                    <span>
                      × {item.quantity}
                    </span>

                    <span>
                      ₹
                      {Number(
                        item.subtotal
                      ).toFixed(2)}
                    </span>
                  </div>
                ))}
              </div>
            )}
          </article>
        ))}
      </section>
    </main>
  );
}

export default Orders;