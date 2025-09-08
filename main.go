/// ATM PROGRAM IN GOLANG

// input from user 
// 1 -> ADD CASH
// 2 -> WITHDRAW CASH
// 3 -> CHECK BALANCE
// 4  -> EXIT

package main
import "fmt"

// MAIN VARIABLE 
 var AMOUNT int = 0 ;
 var CASH int   = 0 ;



// 1 -> ADD CASH FUNCTION
func ADD(){
	fmt.Println("ENTER CASH")
	fmt.Scanln(&CASH)
	AMOUNT+=CASH;
	fmt.Println("Cash added Successfully !!!\n\n\n\n\n");

}

// 2 -> WITHDRAW CASH
func WITHDRAW(){
	fmt.Println("ENTER CASH")
	fmt.Scanln(&CASH)
	AMOUNT-=CASH;
	fmt.Println("Cash Withdrawed Successfully !!!\n\n\n\n\n");
}

// 3 -> CHECK BALANCE
func CHECK(){
	fmt.Println("Your Current Balance is:",AMOUNT,"\n\n\n\n\n");
}




func main() {

	input := 0;

	for i:=0;true;i++{
		fmt.Println(`Enter Choice:
		1-> ADD CASH
		2-> WITHDRAW CASH
		3-> CHECK CASH
		4-> EXIT `);
		fmt.Scanln(&input);
		if input==1 {ADD()}
		if input==2 {WITHDRAW()}
		if input==3 {CHECK()}
		if input==4 {break}
	
   }
}
