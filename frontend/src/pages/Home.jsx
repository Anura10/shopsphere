import { Link } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

function Home() {
  const { user } = useAuth();

  return (
    <main className="home-page">
      <section className="hero-section">
        <div className="hero-content">
          <p className="hero-label">
            WELCOME TO SHOPSPHERE
          </p>

          <h1>
            Everything you need,
            <br />
            all in one place.
          </h1>

          <p className="hero-description">
            Discover quality products, manage your
            cart, and enjoy a simple and secure shopping
            experience.
          </p>

          <div className="hero-actions">
            <Link
              to="/products"
              className="primary-button"
            >
              Explore Products
            </Link>

            <Link
              to="/cart"
              className="secondary-button"
            >
              View Cart
            </Link>
          </div>
        </div>
      </section>

      <section className="categories-section">
        <div className="section-header">
          <p className="section-label">
            SHOP BY CATEGORY
          </p>

          <h2>
            Find what you're looking for
          </h2>
        </div>

        <div className="category-grid">
          <Link
            to="/products?category=electronics"
            className="category-card"
          >
            <span className="category-icon">
              ⚡
            </span>

            <h3>Electronics</h3>

            <p>
              Smart devices and modern technology.
            </p>
          </Link>

          <Link
            to="/products?category=fashion"
            className="category-card"
          >
            <span className="category-icon">
              ◈
            </span>

            <h3>Fashion</h3>

            <p>
              Explore products for every style.
            </p>
          </Link>

          <Link
            to="/products?category=home"
            className="category-card"
          >
            <span className="category-icon">
              ⌂
            </span>

            <h3>Home</h3>

            <p>
              Products designed for everyday life.
            </p>
          </Link>

          <Link
            to="/products"
            className="category-card"
          >
            <span className="category-icon">
              +
            </span>

            <h3>View All</h3>

            <p>
              Browse the complete ShopSphere catalog.
            </p>
          </Link>
        </div>
      </section>

      <section className="welcome-section">
        <div>
          <p className="section-label">
            YOUR SHOPSPHERE
          </p>

          <h2>
            {user?.name
              ? `Welcome back, ${user.name}`
              : "Welcome to ShopSphere"}
          </h2>

          <p>
            Your products, cart and orders are
            connected to the ShopSphere backend.
          </p>
        </div>

        <Link
          to="/products"
          className="primary-button"
        >
          Start Shopping
        </Link>
      </section>
    </main>
  );
}

export default Home;