export default function Home() {
  const highlights = [
    {
      title: 'SMS automation engine',
      text: 'Run personalized SMS sequences with templates, triggers, and timed follow-ups.',
    },
    {
      title: 'CallFlow-style follow-up',
      text: 'For phone-first businesses, send instant missed-call responses and reminder messages.',
    },
    {
      title: 'Shareable profile links',
      text: 'Every business gets a clean public page to convert ad traffic into actions.',
    },
  ];

  const steps = [
    'Set your business profile, automation rules, and templates.',
    'AdFlow captures inbound leads and missed-call events from your campaigns.',
    'Leads receive automated SMS and landing-link follow-up instantly.',
    'They can message, map, or contact your business from one page.',
  ];

  return (
    <main className="container site-home">
      <section className="site-hero">
        <div>
          <p className="label">AdFlow</p>
          <h1>Convert ad traffic into real customers.</h1>
          <p className="site-subtitle">
            AdFlow automates lead engagement with SMS automation and CallFlow-style follow-up, then routes prospects to a focused landing page.
          </p>
          <div className="site-hero-actions">
            <a className="action" href="https://adflow.up.railway.app/admin" target="_blank" rel="noreferrer">
              Open Admin Console
            </a>
            <a className="action action-ghost" href="#how-it-works">
              How It Works
            </a>
          </div>
        </div>

        <div className="site-stat-grid">
          <div className="site-stat-card">
            <p className="site-stat-value">24/7</p>
            <p className="site-stat-label">Automated lead engagement</p>
          </div>
          <div className="site-stat-card">
            <p className="site-stat-value">SMS Automation</p>
            <p className="site-stat-label">Triggers, templates, and sequences</p>
          </div>
          <div className="site-stat-card">
            <p className="site-stat-value">Public page</p>
            <p className="site-stat-label">One shareable business URL</p>
          </div>
        </div>
      </section>

      <section className="section">
        <p className="label">Why AdFlow</p>
        <h2 className="site-section-title">Everything needed for fast lead response</h2>
        <div className="site-grid">
          {highlights.map((item) => (
            <article key={item.title} className="site-card">
              <h3>{item.title}</h3>
              <p>{item.text}</p>
            </article>
          ))}
        </div>
      </section>

      <section id="how-it-works" className="section">
        <p className="label">Workflow</p>
        <h2 className="site-section-title">How your landing flow works</h2>
        <ol className="site-steps">
          {steps.map((step) => (
            <li key={step}>{step}</li>
          ))}
        </ol>
      </section>

      <section className="section site-cta">
        <h2 className="site-section-title">Customer pages are published at:</h2>
        <p className="site-url">https://adflowapp.vercel.app/&lt;customer-id&gt;</p>
        <p className="description">
          Use your assigned customer ID URL to open a specific business landing page.
        </p>
      </section>

      <section className="section">
        <p className="label">Advertising</p>
        <h2 className="site-section-title">Contact for Ads and Partnerships</h2>
        <div className="site-hero-actions">
          <a className="action" href="mailto:mahesh.moholkar.dev@gmail.com">
            mahesh.moholkar.dev@gmail.com
          </a>
          <a className="action action-ghost" href="tel:+919579047391">
            +91 95790 47391
          </a>
        </div>
      </section>
    </main>
  );
}
