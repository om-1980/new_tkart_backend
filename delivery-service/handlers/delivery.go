package handlers

import (
    "database/sql"
    "net/http"
    "delivery-service/models"
    "delivery-service/utils"

    "github.com/labstack/echo/v4"
)

func GetAssignedDeliveries(db *sql.DB) echo.HandlerFunc {
    return func(c echo.Context) error {
        email := c.Get("email").(string)

        rows, err := db.Query(`SELECT * FROM deliveries WHERE delivery_person=$1`, email)
        if err != nil {
            return utils.InternalServerError(c, err)
        }
        defer rows.Close()

        deliveries := []models.Delivery{}
        for rows.Next() {
            var d models.Delivery
            if err := rows.Scan(&d.DeliveryID, &d.OrderID, &d.DeliveryPerson, &d.Status, &d.IsCOD, &d.CODAmount, &d.IsReturn); err != nil {
                return utils.InternalServerError(c, err)
            }
            deliveries = append(deliveries, d)
        }

        return c.JSON(http.StatusOK, deliveries)
    }
}

func UpdateDeliveryStatus(db *sql.DB) echo.HandlerFunc {
    return func(c echo.Context) error {
        var d models.Delivery
        if err := c.Bind(&d); err != nil {
            return utils.BadRequest(c, "Invalid payload")
        }

        _, err := db.Exec(`UPDATE deliveries SET status=$1 WHERE delivery_id=$2`, d.Status, d.DeliveryID)
        if err != nil {
            return utils.InternalServerError(c, err)
        }

        return c.JSON(http.StatusOK, echo.Map{"message": "Status updated"})
    }
}

func RaiseReturnRequest(db *sql.DB) echo.HandlerFunc {
    return func(c echo.Context) error {
        id := c.Param("delivery_id")

        _, err := db.Exec(`UPDATE deliveries SET status='return_initiated', is_return=true WHERE delivery_id=$1`, id)
        if err != nil {
            return utils.InternalServerError(c, err)
        }

        return c.JSON(http.StatusOK, echo.Map{"message": "Return initiated"})
    }
}
