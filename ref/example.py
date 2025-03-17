const CLIENT_ID = "58609595-1742043842907@advertising.performance.ozon.ru"; // Укажите ваш client_id
const CLIENT_SECRET = "ikWg0NgOw5xWYMljx8AvjsWHMm6keI-C9Qc9QMATTFzPdtbbz_kProDF10i7pizmUWZY9vniEX8U4csXUg"; // Укажите ваш client_secret
const TOKEN_URL = "https://api-performance.ozon.ru/api/client/token";
const STATISTICS_URL = "https://api-performance.ozon.ru/api/client/statistics";

let accessToken = ""; // Переменная для хранения токена
let tokenExpiresAt = 0; // Время истечения токена

// Получение токена
function getAccessToken() {
  const now = Math.floor(Date.now() / 1000);

  if (accessToken && now < tokenExpiresAt) {
    Logger.log("Токен ещё действителен.");
    return accessToken;
  }

  const payload = {
    client_id: CLIENT_ID,
    client_secret: CLIENT_SECRET,
    grant_type: "client_credentials"
  };

  const options = {
    method: "post",
    contentType: "application/json",
    payload: JSON.stringify(payload)
  };

  const response = UrlFetchApp.fetch(TOKEN_URL, options);
  const data = JSON.parse(response.getContentText());

  accessToken = data.access_token;
  tokenExpiresAt = Math.floor(Date.now() / 1000) + data.expires_in - 60; // Буфер 60 секунд

  Logger.log("Новый токен получен: " + accessToken);
  return accessToken;
}

// Функция для получения статистики по списку кампаний
function fetchCampaignStatistics() {
  const token = getAccessToken();

  // Укажите список ID кампаний
  const campaigns = [
    "13631009",
    "13631118",
    "13631156",
    "13631365",
    "13706453",
    "13862722",
    "13862730",
    "13865649",
    "13907141",
    "13922577",
    "13922657",
    "13926429",
    "13964241",
    "13964256",
    "13964314",
    "14070254"
  ];

  // Укажите диапазон дат
  const dateFrom = "2025-03-01"; // Начальная дата
  const dateTo = "2025-03-14";   // Конечная дата

  // Формируем тело запроса
  const payload = {
    campaigns: campaigns,
    from: ${dateFrom}T00:00:00Z,
    to: ${dateTo}T23:59:59Z,
    groupBy: "NO_GROUP_BY" // Группировка (можно изменить на нужную)
  };

  const options = {
    method: "post",
    contentType: "application/json",
    headers: {
      Authorization: "Bearer " + token
    },
    payload: JSON.stringify(payload),
    muteHttpExceptions: true
  };

  const response = UrlFetchApp.fetch(STATISTICS_URL, options);
  const responseText = response.getContentText();

  if (!responseText) {
    throw new Error("Сервер вернул пустой ответ.");
  }

  if (response.getResponseCode() !== 200) {
    Logger.log("Код ошибки: " + response.getResponseCode());
    throw new Error("Ошибка: " + responseText);
  }

  const data = JSON.parse(responseText);

  // Проверяем, есть ли данные
  if (!data  !data.result  data.result.length === 0) {
    Logger.log("Нет данных для отображения.");
    return [];
  }

  // Преобразуем данные в массив для Google Таблиц
  const rows = [];
  rows.push(Object.keys(data.result[0])); // Заголовки колонок
  data.result.forEach(item => {
    rows.push(Object.values(item)); // Значения
  });

  return rows;
}

// Экспорт данных в Google Таблицы
function exportCampaignStatisticsToGoogleSheets() {
  const reportData = fetchCampaignStatistics();

  if (reportData.length === 0) {
    Logger.log("Нет данных для выгрузки.");
    return;
  }

  const sheet = SpreadsheetApp.getActiveSpreadsheet().getActiveSheet();
  sheet.clear();

  sheet.getRange(1, 1, reportData.length, reportData[0].length).setValues(reportData);

  Logger.log("Данные успешно выгружены в Google Таблицы!");
}