# Frontend JavaScript Module Structure

## Overview

This is the frontend JavaScript module structure for the medical device manufacturer website. The code is organized into modular components following modern web development best practices.

## Directory Structure

```
frontend/js/
├── main.js                    # Entry point
├── components/                # Page components
│   ├── Header.js              # Header navigation and language switcher
│   ├── Banner.js              # Hero banner section
│   ├── About.js               # About us section
│   ├── Stats.js               # Statistics section
│   ├── Products.js            # Products gallery section
│   ├── Factory.js             # Factory strength section
│   ├── Advantage.js           # Company advantages section
│   ├── Events.js              # Upcoming events section
│   ├── Contact.js             # Contact form section
│   └── Brand.js               # Brand/logo management
├── services/                  # API services
│   └── cmsService.js          # CMS data fetching and management
├── utils/                     # Utility functions
│   └── i18n.js                # Internationalization support
└── config/                    # Configuration files (reserved)
```

## Modules

### 1. Entry Point - main.js

The main entry point that initializes all components and services.

**Key functions:**
- Initializes i18n system
- Initializes header component
- Loads CMS data
- Applies CMS content to all sections

**Usage:**
```html
<script type="module" src="/frontend/js/main.js"></script>
```

### 2. Components

#### Header.js
- **Responsibilities:** Language switcher, mobile menu toggle, navigation link handling
- **Exports:** `initHeader()`

#### Banner.js
- **Responsibilities:** Loads and renders banner content from CMS
- **Exports:** `loadBannerContent()`

#### About.js
- **Responsibilities:** Loads and renders about section content
- **Exports:** `loadAboutContent()`

#### Stats.js
- **Responsibilities:** Loads and renders statistics data
- **Exports:** `loadStatsContent()`

#### Products.js
- **Responsibilities:** Loads and renders product gallery
- **Exports:** `loadProductsContent()`

#### Factory.js
- **Responsibilities:** Loads and renders factory strength section
- **Exports:** `loadFactoryContent()`

#### Advantage.js
- **Responsibilities:** Loads and renders company advantages
- **Exports:** `loadAdvantageContent()`

#### Events.js
- **Responsibilities:** Loads and renders upcoming events
- **Exports:** `loadEventContent()`

#### Contact.js
- **Responsibilities:** Loads and renders contact information
- **Exports:** `loadContactContent()`

#### Brand.js
- **Responsibilities:** Loads and renders brand/logo information
- **Exports:** `loadBrandContent()`

### 3. Services

#### cmsService.js
- **Responsibilities:** Fetch CMS data from API, manage CMS state, provide language-specific field access
- **Exports:** 
  - `loadCMSData()` - Fetch data from API
  - `getCMSData()` - Get current CMS data
  - `getLangSpecificField(obj, baseField)` - Get language-specific field value
  - `getCurrentLangFromStorage()` - Get current language from localStorage

### 4. Utilities

#### i18n.js
- **Responsibilities:** Internationalization support, language switching, translation management
- **Exports:**
  - `initI18n()` - Initialize i18n system
  - `setLang(lang)` - Set current language
  - `getCurrentLang()` - Get current language
  - `translate(key)` - Translate a key
  - `applyTranslations()` - Apply translations to DOM
  - `loadTranslationsFromAPI()` - Load dynamic translations

## Dependencies

```
main.js
├── utils/i18n.js
├── components/Header.js
│   └── utils/i18n.js
├── components/Banner.js
│   └── services/cmsService.js
├── components/About.js
│   └── services/cmsService.js
├── components/Stats.js
│   └── services/cmsService.js
├── components/Products.js
│   └── services/cmsService.js
├── components/Factory.js
│   └── services/cmsService.js
├── components/Advantage.js
│   └── services/cmsService.js
├── components/Events.js
│   └── services/cmsService.js
├── components/Contact.js
│   └── services/cmsService.js
└── components/Brand.js
    └── services/cmsService.js
```

## Language Support

- English (en)
- Chinese (zh)

Language is stored in localStorage as `preferredLang`.

## Data Flow

1. Page loads → `main.js` initializes
2. `initI18n()` sets up language system
3. `initHeader()` sets up navigation
4. `loadCMSData()` fetches data from API
5. `applyCMSContent()` applies content to all components
6. Language switch triggers re-application of translations and CMS content