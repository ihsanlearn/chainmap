package core

const DefaultXSLT = `<?xml version="1.0" encoding="UTF-8"?>
<xsl:stylesheet version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
  <xsl:output method="html" doctype-public="XSLT-compat" encoding="utf-8" indent="yes"/>
  <xsl:template match="/">
    <html lang="en">
      <head>
        <title>Chainmap Scan Report</title>
        <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
        <style>
          :root {
            --primary: #0f172a;
            --secondary: #334155;
            --accent: #3b82f6;
            --bg: #f8fafc;
            --card-bg: #ffffff;
            --text: #1e293b;
            --text-light: #64748b;
            --success: #10b981;
            --warning: #f59e0b;
            --danger: #ef4444;
            --border: #e2e8f0;
          }
          body {
            font-family: 'Inter', -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
            background-color: var(--bg);
            color: var(--text);
            margin: 0;
            padding: 0;
            line-height: 1.5;
          }
          .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 2rem;
          }
          header {
            background-color: var(--card-bg);
            padding: 1.5rem 2rem;
            border-bottom: 1px solid var(--border);
            margin-bottom: 2rem;
            box-shadow: 0 1px 3px rgba(0,0,0,0.05);
            display: flex;
            justify-content: space-between;
            align-items: center;
          }
          .logo {
            font-size: 1.5rem;
            font-weight: 700;
            color: var(--primary);
            text-decoration: none;
            display: flex;
            align-items: center;
            gap: 0.5rem;
          }
          .badge {
            display: inline-block;
            padding: 0.25rem 0.75rem;
            border-radius: 9999px;
            font-size: 0.75rem;
            font-weight: 600;
            text-transform: uppercase;
          }
          .badge-success { background-color: #d1fae5; color: #065f46; }
          .badge-warning { background-color: #fef3c7; color: #92400e; }
          .badge-danger { background-color: #fee2e2; color: #991b1b; }
          
          .dashboard-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 1.5rem;
            margin-bottom: 2rem;
          }
          .card {
            background: var(--card-bg);
            border-radius: 0.75rem;
            padding: 1.5rem;
            box-shadow: 0 1px 2px rgba(0,0,0,0.05);
            border: 1px solid var(--border);
          }
          .stat-title {
            color: var(--text-light);
            font-size: 0.875rem;
            font-weight: 500;
            margin-bottom: 0.5rem;
          }
          .stat-value {
            font-size: 2rem;
            font-weight: 700;
            color: var(--primary);
          }
          
          .host-card {
            background: var(--card-bg);
            border-radius: 0.75rem;
            box-shadow: 0 1px 2px rgba(0,0,0,0.05);
            border: 1px solid var(--border);
            margin-bottom: 1.5rem;
            overflow: hidden;
          }
          .host-header {
            padding: 1.25rem;
            border-bottom: 1px solid var(--border);
            display: flex;
            justify-content: space-between;
            align-items: center;
            background-color: #f8fafc;
          }
          .host-title {
            font-size: 1.125rem;
            font-weight: 600;
            color: var(--primary);
            display: flex;
            align-items: center;
            gap: 0.75rem;
          }
          .host-meta {
            font-size: 0.875rem;
            color: var(--text-light);
          }
          
          table {
            width: 100%;
            border-collapse: collapse;
          }
          th {
            background-color: #f1f5f9;
            text-align: left;
            padding: 0.75rem 1.25rem;
            font-size: 0.75rem;
            font-weight: 600;
            text-transform: uppercase;
            color: var(--text-light);
            border-bottom: 1px solid var(--border);
          }
          td {
            padding: 1rem 1.25rem;
            border-bottom: 1px solid var(--border);
            font-size: 0.875rem;
            vertical-align: top;
          }
          tr:last-child td { border-bottom: none; }
          .port-open { color: var(--success); font-weight: 600; }
          .port-closed { color: var(--danger); }
          .port-filtered { color: var(--warning); }
          
          .script-output {
            margin-top: 0.5rem;
            background-color: #f8fafc;
            border: 1px solid #e2e8f0;
            border-radius: 0.375rem;
            padding: 0.75rem;
            font-family: monospace;
            font-size: 0.8rem;
            white-space: pre-wrap;
            color: #334155;
          }
          .script-id {
            color: #475569;
            font-weight: 600;
            margin-bottom: 0.25rem;
            display: block;
          }
          
        </style>
      </head>
      <body>
        <header>
          <div class="logo">
            <span>⚡ Chainmap Report</span>
          </div>
          <div class="host-meta">
            Generated: <xsl:value-of select="/NmapRun/@startstr"/>
          </div>
        </header>

        <div class="container">
          <!-- Dashboard -->
          <div class="dashboard-grid">
            <div class="card">
              <div class="stat-title">Total Targets</div>
              <div class="stat-value"><xsl:value-of select="/NmapRun/runstats/hosts/@total"/></div>
            </div>
            <div class="card">
              <div class="stat-title">Hosts Up</div>
              <div class="stat-value" style="color: var(--success)">
                <xsl:value-of select="/NmapRun/runstats/hosts/@up"/>
              </div>
            </div>
            <div class="card">
              <div class="stat-title">Scan Duration</div>
              <div class="stat-value"><xsl:value-of select="/NmapRun/runstats/finished/@elapsed"/>s</div>
            </div>
          </div>

          <!-- Hosts List -->
          <xsl:for-each select="/NmapRun/host">
            <div class="host-card">
              <div class="host-header">
                <div class="host-title">
                  <xsl:value-of select="address/@addr"/>
                  <xsl:if test="hostnames/hostname/@name">
                     <span style="color: var(--text-light); font-weight: 400; font-size: 0.9em;">
                       (<xsl:value-of select="hostnames/hostname/@name"/>)
                     </span>
                  </xsl:if>
                </div>
                <div class="badge badge-success">UP</div>
              </div>
              
              <xsl:choose>
                <xsl:when test="ports/port">
                  <table>
                    <thead>
                      <tr>
                        <th style="width: 15%">Port</th>
                        <th style="width: 15%">State</th>
                        <th style="width: 30%">Service</th>
                        <th style="width: 40%">Version / Scripts</th>
                      </tr>
                    </thead>
                    <tbody>
                      <xsl:for-each select="ports/port">
                        <tr>
                          <td><strong><xsl:value-of select="@portid"/></strong>/<xsl:value-of select="@protocol"/></td>
                          <td>
                            <xsl:choose>
                              <xsl:when test="state/@state = 'open'">
                                <span class="badge badge-success">OPEN</span>
                              </xsl:when>
                              <xsl:when test="state/@state = 'filtered'">
                                <span class="badge badge-warning">FILTERED</span>
                              </xsl:when>
                              <xsl:otherwise>
                                <span class="badge badge-danger"><xsl:value-of select="state/@state"/></span>
                              </xsl:otherwise>
                            </xsl:choose>
                            <xsl:if test="state/@reason">
                              <div style="font-size: 0.75rem; color: var(--text-light); margin-top: 0.25rem;">
                                <xsl:value-of select="state/@reason"/>
                                <xsl:if test="state/@reason_ttl">
                                  (ttl <xsl:value-of select="state/@reason_ttl"/>)
                                </xsl:if>
                              </div>
                            </xsl:if>
                          </td>
                          <td style="color: var(--accent); font-weight: 500;">
                            <xsl:value-of select="service/@name"/>
                          </td>
                          <td style="color: var(--text-light);">
                            <div>
                              <xsl:value-of select="service/@product"/>
                              <xsl:if test="service/@version">
                                v<xsl:value-of select="service/@version"/>
                              </xsl:if>
                              <xsl:if test="service/@extrainfo">
                                <span style="font-size: 0.8em; margin-left: 4px; opacity: 0.8">
                                  (<xsl:value-of select="service/@extrainfo"/>)
                                </span>
                              </xsl:if>
                            </div>

                            <!-- Script Outputs -->
                            <xsl:for-each select="script">
                              <div class="script-output">
                                <span class="script-id"><xsl:value-of select="@id"/></span>
                                <xsl:value-of select="@output"/>
                              </div>
                            </xsl:for-each>
                          </td>
                        </tr>
                      </xsl:for-each>
                    </tbody>
                  </table>
                </xsl:when>
                <xsl:otherwise>
                  <div style="padding: 1.5rem; color: var(--text-light); font-style: italic;">
                    No open ports found.
                  </div>
                </xsl:otherwise>
              </xsl:choose>
            </div>
          </xsl:for-each>
        </div>
      </body>
    </html>
  </xsl:template>
</xsl:stylesheet>
`
