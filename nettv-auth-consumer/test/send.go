package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	_ "github.com/godror/godror"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	// db, err := sql.Open("godror", `user="ebill" password="dVri4sd04Be4ZpU_OwwwIXE4VA" connectString="rac-scan.wlink.com.np:1521/raddb_apinettv"`)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// defer db.Close()

	// rows, err := db.Query("SELECT * FROM customer_ebill ce, customer_info  ci WHERE ce.user_name = ci.machine_name AND ce.user_name = 'ranjanmanandhar_home'")

	// cols, _ := rows.Columns()
	// if err != nil {
	// 	fmt.Println("Error running query")
	// 	fmt.Println(err)
	// 	return
	// }
	// defer rows.Close()

	// // var username driver.Rows
	// m := make(map[string]interface{})

	// for rows.Next() {
	// 	columns := make([]interface{}, len(cols))
	// 	columnPointers := make([]interface{}, len(cols))
	// 	for i, _ := range columns {
	// 		columnPointers[i] = &columns[i]
	// 	}

	// 	// Scan the result into the column pointers...
	// 	if err := rows.Scan(columnPointers...); err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	// Create our map, and retrieve the value for each column from the pointers slice,
	// 	// storing it in the map with the name of the column as the key.
	// 	// m := make(map[string]interface{})
	// 	for i, colName := range cols {
	// 		val := columnPointers[i].(*interface{})
	// 		m[colName] = *val
	// 	}

	// 	// Outputs: map[columnName:value columnName2:value2 columnName3:value3 ...]
	// 	// fmt.Print(m)

	// 	// fmt.Printf(fmt.Sprintf("The username is: %+v, mac_id is: %s\n", username, stb_box_id))
	// }
	// // fmt.Println(m)
	// rows, _ = db.Query("SELECT stb_box_id FROM wlinkiptv_usermap WHERE username = 'ranjanmanandhar_home'")

	// cols, _ = rows.Columns()
	// data := make(map[string]string)

	// if rows.Next() {
	// 	columns := make([]string, len(cols))
	// 	columnPointers := make([]interface{}, len(cols))
	// 	for i, _ := range columns {
	// 		columnPointers[i] = &columns[i]
	// 	}

	// 	rows.Scan(columnPointers...)

	// 	for i, colName := range cols {
	// 		data[colName] = columns[i]
	// 	}
	// }
	// fmt.Println(data)
	// var newStr string
	// var rset1 driver.Rows
	// sqlstmt := "DECLARE v_pay_plan_id VARCHAR2(32); v_pay_plan_name VARCHAR2(32); v_plan_category_id VARCHAR2(32); v_plan_category_name VARCHAR2(32); v_acct_check_by VARCHAR2(32); v_up_bw VARCHAR2(32); v_down_bw VARCHAR2(32); v_min_up_bw VARCHAR2(32); v_min_down_bw VARCHAR2(32); v_grace_status VARCHAR2(32); v_expiry_date VARCHAR2(32); v_disable VARCHAR2(32); v_online_date VARCHAR2(32); v_days_remaining VARCHAR2(32); v_minutes_remaining VARCHAR2(32); v_volume_remaining VARCHAR2(32); v_extend_status VARCHAR2(32); v_extend_date VARCHAR2(32); v_extend_days VARCHAR2(32); v_extend_hrs VARCHAR2(32); v_extend_volume VARCHAR2(32); v_sme_plan VARCHAR2(32); v_support_zone VARCHAR2(32); v_supportzone_id VARCHAR2(32); BEGIN ebill.account_status_v2( :user_name, v_pay_plan_id, v_pay_plan_name, v_plan_category_id, v_plan_category_name, v_acct_check_by, v_up_bw, v_down_bw, v_min_up_bw, v_min_down_bw, v_grace_status, v_expiry_date, v_disable, v_online_date, v_days_remaining, v_minutes_remaining, v_volume_remaining, v_extend_status, v_extend_date, v_extend_days, v_extend_hrs, v_extend_volume, v_sme_plan, v_support_zone, v_supportzone_id ); --    dbms_output.put_line('user_name: ' || v_user_name); dbms_output.put_line('pay_plan_id: ' || v_pay_plan_id); dbms_output.put_line('pay_plan_name: ' || v_pay_plan_name); dbms_output.put_line('plan_category_id: ' || v_plan_category_id); dbms_output.put_line('plan_category_name: ' || v_plan_category_name); dbms_output.put_line('acct_check_by: ' || v_acct_check_by); dbms_output.put_line('up_bw: ' || v_up_bw); dbms_output.put_line('down_bw: ' || v_down_bw); dbms_output.put_line('min_up_bw: ' || v_min_up_bw); dbms_output.put_line('min_down_bw: ' || v_min_down_bw); dbms_output.put_line('grace_status: ' || v_grace_status); dbms_output.put_line('expiry_date: ' || v_expiry_date); dbms_output.put_line('disable: ' || v_disable); dbms_output.put_line('online_date: ' || v_online_date); dbms_output.put_line('days_remaining: ' || v_days_remaining); dbms_output.put_line('minutes_remaining: ' || v_minutes_remaining); dbms_output.put_line('volume_remaining: ' || v_volume_remaining); dbms_output.put_line('extend_status: ' || v_extend_status); dbms_output.put_line('extend_date: ' || v_extend_date); dbms_output.put_line('extend_days: ' || v_extend_days); dbms_output.put_line('extend_hrs: ' || v_extend_hrs); dbms_output.put_line('extend_volume: ' || v_extend_volume); dbms_output.put_line('sme_plan: ' || v_sme_plan); dbms_output.put_line('support_zone: ' || v_support_zone); dbms_output.put_line('supportzone_id: ' || v_supportzone_id); END;"
	// // stmt, _ := db.Exec(sqlstmt, sql.Named("username", "ranjanmanandhar_home"), sql.Named("data", sql.Out{Dest: &newStr}))

	// // const query = `BEGIN Package.StoredProcA(123, :1, :2); END;`
	// ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// stmt, _ := db.PrepareContext(ctx, sqlstmt)

	// defer stmt.Close()
	// if _, err := stmt.ExecContext(ctx, sql.Out{Dest: &rset1}); err != nil {
	// 	log.Printf("Error running %q: %+v", sqlstmt, err)
	// 	return
	// }
	// defer rset1.Close()

	// cols1 := rset1.(driver.RowsColumnTypeScanType).Columns()
	// dests1 := make([]driver.Value, len(cols1))
	// for {
	// 	if err := rset1.Next(dests1); err != nil {
	// 		if err == io.EOF {
	// 			break
	// 		}
	// 		rset1.Close()
	// 		// return err
	// 	}
	// 	fmt.Println(dests1)
	// }
	var rabbitmq = flag.String("rabbitmq", "localhost:5671", "Rabbitmq URL")
	var mongoURL = flag.String("mongo", "localhost:27017", "Monog URL")
	flag.Parse()
	conn, err := amqp.Dial(fmt.Sprintf("amqp://guest:guest@%s", *rabbitmq))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	// q, err := ch.QueueDeclare(
	// 	"test", // name
	// 	false,  // durable
	// 	false,  // delete when unused
	// 	false,  // exclusive
	// 	false,  // no-wait
	// 	nil,    // arguments
	// )
	// failOnError(err, "Failed to declare a queue")

	// body := `{"username":"pradi",
	// "data":{
	// 	"name":"ranjan"
	// }}`
	body := `{
	        "username" : "lila_fblri",
	        "stb_box_ids" : "00226DB3539G",
	        "new_netv_map" : true,
	        "data_type" : "create",
	        "type" : "nettv"
	}`
	// body := `{
	//         "username" : "H594-K01-01",
	//         "stb_box_ids" : "00:22:6D:91:FF:6A",
	// 		"address" : "Nayabazar",
	// 		"Email": "eclothing2005@yahoo.com",
	// 		"Phone": "2147483647",
	// 		"Mobile": "2147483647",
	//         "new_netv_map" : true,
	// 		"company_name": "Hotel Kesu Home Pvt. Ltd.",
	//         "type" : "corporate"
	// 	}`
	// body := `{
	// 	"type": "customer",
	// 	"account_type": "PERSONAL",
	// 	"disable": "N",
	// 	"installed_by": "roshan.twanabasu",
	// 	"longitude": "",
	// 	"lattitude": "",
	// 	"machine_name": "test7777",
	// 	"marketed_by": "roshan.twanabasu",
	// 	"mem_expiry_date": "2023-04-18",
	// 	"mem_start_date": "2022-04-18",
	// 	"pay_plan": "1888",
	// 	"per_address": "Jawlakhel, Lalitpur",
	// 	"per_client_name": "Test User",
	// 	"per_cont_email_primary": "test@gmail.com",
	// 	"per_cont_email_secondary": "secondary@gmail.com",
	// 	"per_cont_mobile": "9811111111",
	// 	"per_cont_mobile_secondary": "9711111111",
	// 	"per_cont_phone1": "1-5050507",
	// 	"per_cont_phone2": "1-5050507",
	// 	"per_district": "Jawlakhel",
	// 	"per_fax": "1-5050507",
	// 	"per_house_no": "121",
	// 	"per_municipality": "Lalitpur Municipality",
	// 	"per_pobox": "",
	// 	"per_street": "Jawlakhel Road",
	// 	"per_url": "",
	// 	"per_ward_no": "05",
	// 	"plan_category_id": "1888123",
	// 	"referred_by": "Roshan",
	// 	"vatno": "1231231",
	// 	"tt_ni_ticketid": "",
	// 	"broutermake": "",
	// 	"broutermodel": "",
	// 	"brouterversion": "",
	// 	"conn_setup_area": "1-5050507",
	// 	"org_name": "New Organization",
	// 	"org_address": "Jawlakhel, Lalitpur",
	// 	"org_cont_email_primary": "primaryorg@gmail.com",
	// 	"org_cont_email_secondary": "secondaryorg@gmail.com",
	// 	"org_cont_mobile": "9811111111",
	// 	"org_cont_mobile_secondary": "9711111111",
	// 	"org_cont_name": "neworganization",
	// 	"org_cont_phone1": "1-5512312",
	// 	"org_cont_phone2": "1-5512312",
	// 	"org_district": "Lalitpur",
	// 	"org_fax": "1-5512312",
	// 	"org_house_no": "120",
	// 	"org_municipality": "Lalitpur Municipality",
	// 	"org_pobox": "",
	// 	"org_street": "Jawlakhel Road",
	// 	"org_type": "INGO",
	// 	"org_url": "neworg.com.np",
	// 	"org_ward_no": "05",
	// 	"finan_cont_name": "Roshan",
	// 	"finan_cont_address": "Jawlakhel, Lalitpur",
	// 	"finan_cont_email_primary": "primaryfin@gmail.com",
	// 	"finan_cont_email_secondary": "secondaryfin@gmail.com",
	// 	"finan_cont_mobile": "9811111111",
	// 	"finan_cont_phone1": "1-5050507",
	// 	"finan_cont_phone2": "1-5050507",
	// 	"finan_fax": "1-5050507",
	// 	"finan_pobox": "",
	// 	"tech_cont_name": "Roshan",
	// 	"tech_cont_address": "Jawlakhel, Lalitpur",
	// 	"tech_cont_email_primary": "tech@gmail.com",
	// 	"tech_cont_email_secondary": "tech@gmail.com",
	// 	"tech_cont_mobile": "9811111111",
	// 	"tech_cont_phone1": "1-5050507",
	// 	"tech_cont_phone2": "1-5050507",
	// 	"tech_fax": "1-5050507",
	// 	"tech_pobox": "",
	// 	"allow_nos": "",
	// 	"caller_id_only": "",
	// 	"create_date": "2022-04-18",
	// 	"created_by": "roshan.twanabasu",
	// 	"deny_nos": "",
	// 	"dur_min_left": "",
	// 	"expiry_date": "2023-04-18",
	// 	"extra_days": "1",
	// 	"extra_hrs": "1",
	// 	"extra_volume": "1",
	// 	"ippoe_mac": "",
	// 	"ippoe_convert_date": "2022-04-18",
	// 	"last_admin_use": "",
	// 	"last_dialup_use": "",
	// 	"last_pop_use": "",
	// 	"last_session_id": "",
	// 	"latest_payment_amt": 1000,
	// 	"mail_quota_size": "",
	// 	"main_user_name": "test7777",
	// 	"no_disable": "",
	// 	"no_rem": "",
	// 	"online_date": "2022-04-18",
	// 	"online_dur_min": "",
	// 	"online_dur_today": "",
	// 	"pass_word_admin": "",
	// 	"pass_word_dialup": "",
	// 	"pass_word_pop": "",
	// 	"nt_password": "",
	// 	"rem_days_date": "",
	// 	"rem_days_left": "",
	// 	"rem_min_date": "",
	// 	"rem_min_left": "",
	// 	"remarks": "Test customer create",
	// 	"renew_date": "2023-04-18",
	// 	"renew_dur_min": "",
	// 	"renew_dur_month": "",
	// 	"status": "",
	// 	"support_zone": "MANBHAWAN",
	// 	"support_zone_id": "1",
	// 	"u_daily_quota": "",
	// 	"u_idle_timeout": "",
	// 	"u_logins": "",
	// 	"u_monthly_qouta": "",
	// 	"u_session_timeout": "",
	// 	"u_time_value": "",
	// 	"use_allow": "",
	// 	"use_deny": "",
	// 	"user_name": "pradimna",
	// 	"valid_until": "",
	// 	"verified_date": "",
	// 	"verified": "",
	// 	"duration_left": 1000,
	// 	"latest_payment_amount": 1000,
	// 	"notify_hours": 2
	//   }
	// `
	// jsonStr, err := json.Marshal(body)
	err = ch.Publish(
		"nettv_auth_consumer", // exchange
		"",                    // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)

	connectString := fmt.Sprintf("mongodb://%s", *mongoURL)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectString))

	if err != nil {
		panic(err)
	}

	log.Printf("connect to db successful")

	collection := client.Database("nettv_auth_consumer").Collection("nettv_auth_consumer")
	count, err := collection.CountDocuments(ctx, bson.M{"username": "ranjan"})
	fmt.Println(count)

}
