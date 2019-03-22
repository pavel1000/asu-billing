package main

import (
	"fmt"
)

func getUsersByType(t string) ([]User, error) {
	rows, err := db.Query(`SELECT Users.ID, Users.Name, Users.Login, Users.Money, Users.Active, Users.Phone,
		Users.Comment, Users.Payments_ends, In_IPs.IP, Ext_IPs.IP, Tariffs.ID, Tariffs.Name, Tariffs.Price
	FROM (((Users
		INNER JOIN In_IPs ON Users.In_IP_ID = In_IPs.ID)
		INNER JOIN Ext_IPs ON Users.Ext_IP_ID = Ext_IPs.ID)
		INNER JOIN Tariffs ON Users.Tariff_ID = Tariffs.ID)`)
	if err != nil {
		return nil, fmt.Errorf("could not do query: %v", err)
	}
	defer rows.Close()

	var user User
	users := make([]User, 0)
	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Name, &user.Login, &user.Money, &user.Active, &user.Phone, &user.Comment,
			&user.PaymentsEnds, &user.InIP, &user.ExtIP, &user.Tariff.ID, &user.Tariff.Name, &user.Tariff.Price)
		if err != nil {
			return nil, fmt.Errorf("could not scan from row: %v", err)
		}
		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("something happened with rows: %v", err)
	}

	return users, nil
}

func getUserByID(id int) (User, error) {
	var user User
	err := db.QueryRow(`SELECT Users.Name, Users.Login, Users.Money, Users.Active, Users.Phone,
	 	Users.Comment, Users.Payments_ends, In_IPs.IP, Ext_IPs.IP, Tariffs.ID, Tariffs.Name, Tariffs.Price
	FROM (((Users
		INNER JOIN In_IPs ON Users.In_IP_ID = In_IPs.ID)
		INNER JOIN Ext_IPs ON Users.Ext_IP_ID = Ext_IPs.ID)
		INNER JOIN Tariffs ON Users.Tariff_ID = Tariffs.ID)
	WHERE Users.ID = ?`, id).Scan(&user.Name, &user.Login, &user.Money, &user.Active, &user.Phone,
		&user.Comment, &user.PaymentsEnds, &user.InIP, &user.ExtIP, &user.Tariff.ID, &user.Tariff.Name, &user.Tariff.Price)
	if err != nil {
		return user, fmt.Errorf("could not do queryRow: %v", err)
	}

	rows, err := db.Query(`SELECT Amount, Date FROM Payments WHERE User_ID= ?`, id)
	if err != nil {
		return user, fmt.Errorf("could not do query: %v", err)
	}

	var payment Payment
	payments := make([]Payment, 0)
	for rows.Next() {
		err := rows.Scan(&payment.Amount, &payment.Last)
		if err != nil {
			return user, fmt.Errorf("could not scan from row: %v", err)
		}
		payments = append(payments, payment)
	}
	err = rows.Err()
	if err != nil {
		return user, fmt.Errorf("something happened with rows: %v", err)
	}

	user.ID = id
	user.Payments = payments
	return user, nil
}
