package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"time"
)


var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Use all scan mode",
	PreRun: func(cmd *cobra.Command, args []string) {
		CreatFile(Output_result,Path_result)
		PrintScanBanner("all")
	},
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		defer func() {
			Output_endtime(start)
		}()
		allmode()
	},
}

func allmode()  {
	GetHost()
	if pingbefore {
		Hosts = ping_discover()
	}
	ips, err := Parse_IP(Hosts)
	Checkerr(err)
	ports, err := Parse_Port(ps_port)
	Checkerr(err)
	aliveserver:=NewPortScan(ips,ports,Connectall)
	r:=aliveserver.Run()
	getHttptitle(r)
	Printresult(r)
}

func Connectall(ip string, port int) (string, int, error,[]string) {
	var r []string //返回从该端口获取的信息
	conn,err:=Getconn(fmt.Sprintf("%v:%v",ip,port))
	if conn != nil {
		defer conn.Close()
		fmt.Printf(White(fmt.Sprintf("\rFind port %v:%v\r\n", ip, port)))
		switch port {
		case 22:
			fmt.Println(Yellow("\rStart burp ssh : ",ip,":",port))
			name:="root,admin,ssh"
			_,f,_:=ssh_auto("root","Ksdvfjsxc",ip)
			if f{
				Output(fmt.Sprintf("[-]%v Don't allow root login:%v \n","ssh",ip),Yellow)
				name="admin,ssh"
			}
			startburp:=NewBurp(Password,name,Userdict,Passdict,ip,ssh_auto,10)
			relust:=startburp.Run()
			if relust!=""{
				return ip,port,nil,[]string{relust}
			}
			return ip,port,nil,nil
		case 3306:
			fmt.Println(Yellow("\rStart burp mysql : ",ip,":",port))
			_,f,_:=mysql_auth("asdasd","zxczxc",ip)
			if f{
				Output(fmt.Sprintf("[+]%v burp success:%v No authentication\n","mysql",ip),LightGreen)
				return ip,port,nil,[]string{"No authentication"}
			}
			startburp:=NewBurp(Password,"root,mysql",Userdict,Passdict,ip,mysql_auth,100)
			relust:=startburp.Run()
			if relust!=""{
				return ip,port,nil,[]string{relust}
			}
			return ip,port,nil,nil
		case 6379:
			fmt.Println(Yellow("\rStart burp redis : ",ip,":",port))
			_,f,_:=redis_auth("","",ip)
			if f{
				Output(fmt.Sprintf("[+]%v burp success:%v No authentication\n","redis",ip),LightGreen)
				return ip,port,nil,[]string{"No authentication"}
			}
			startburp:=NewBurp(Password,"","",Passdict,ip,redis_auth,100)
			relust:=startburp.Run()
			if relust!=""{
				return ip,port,nil,[]string{relust}
			}
			return ip,port,nil,nil
		case 1433:
			fmt.Println(Yellow("\rStart burp mssql : ",ip,":",port))
			startburp:=NewBurp(Password,"sa,admin,Administrator",Userdict,Passdict,ip,mssql_auth,100)
			relust:=startburp.Run()
			if relust!=""{
				return ip,port,nil,[]string{relust}
			}
			return ip,port,nil,nil
		case 5432:
			fmt.Println(Yellow("\rStart burp postgres : ",ip,":",port))
			startburp:=NewBurp(Password,"postgres",Userdict,Passdict,ip,postgres_auth,100)
			relust:=startburp.Run()
			if relust!=""{
				return ip,port,nil,[]string{relust}
			}
			return ip,port,nil,nil
		case 7890:
			b,s:=Socks5Find(conn)
			if b {
				Output(fmt.Sprintf("%v\t%v:%v \n",s,ip,port),LightGreen)
				r=[]string{s}
				return ip, port, nil,r
			}else {
				return ip,port,nil,nil
			}
		case 10808:
			b,s:=Socks5Find(conn)
			if b {
				Output(fmt.Sprintf("%v\t%v:%v \n",s,ip,port),LightGreen)
				r=[]string{s}
				return ip, port, nil,r
			}else {
				return ip,port,nil,nil
			}
		case 1089:
			b,s:=Socks5Find(conn)
			if b {
				Output(fmt.Sprintf("%v\t%v:%v \n",s,ip,port),LightGreen)
				r=[]string{s}
				return ip, port, nil,r
			}else {
				return ip,port,nil,nil
			}
		case 445:
			_,smbRes:=smbinfo(conn)
			_,_,_,r=Connect17010(ip,port)
			for _,i:=range r{
				smbRes=append(smbRes,i)
			}
			return ip,port,nil,smbRes
		case 135:
			_,oxidres:=oxidIpInfo(conn)
			return ip,port,nil,oxidres
		case 139:
			nbname, _ := netBIOS(ip)
			if nbname.msg != "" {
				return ip, port, nil, []string{nbname.msg}
			}
		case 21:
			if Verbose{
				fmt.Println(Yellow("\rStart burp ftp : ",ip,":",port))
			}
			startburp:=NewBurp(Password,Username,Userdict,Passdict,ip,ftp_auth,burpthread)
			relust:=startburp.Run()
			if relust!=""{
				return ip,port,nil,[]string{relust}
			}
			return ip,port,nil,nil
		case 27017:
			fmt.Println(Yellow("\rStart burp mongodb : ",ip,":",port))
			_,f,_:=mongodb_auth("","",ip)
			if f{
				Output(fmt.Sprintf("[+]%v burp success:%v No authentication\n","mongodb",ip),LightGreen)
				return ip,port,nil,[]string{"No authentication"}
			}
			startburp:=NewBurp(Password,Username,Userdict,Passdict,ip,mongodb_auth,burpthread)
			relust:=startburp.Run()
			if relust!=""{
				return ip,port,nil,[]string{relust}
			}
			return ip,port,nil,nil
		default:
			//getScanTitl(fmt.Sprintf("%v:%v",ip,port))
		}
	}
	return ip, port, err,r
}

func init() {
	rootCmd.AddCommand(allCmd)
	allCmd.Flags().StringVar(&Hostfile,"hostfile","","Set host file")
	allCmd.Flags().BoolVarP(&useicmp,"icmp","i",false,"Icmp packets are sent to check whether the host is alive(need root)")
	allCmd.Flags().StringVarP(&Hosts, "host", "H", "", "Set `hosts`(The format is similar to Nmap) eg:192.168.1.1/24,172.16.95.1-100,127.0.0.1")
	allCmd.Flags().StringVarP(&ps_port, "port", "p", default_port, "Set `port` eg:1-1000,3306,3389")
	allCmd.Flags().BoolVar(&pingbefore, "ping", false, "Ping host discovery before port scanning")
	allCmd.Flags().StringVarP(&Password,"password","P","","Set postgres password")
	allCmd.Flags().StringVarP(&Passdict,"passdict","","","Set postgres passworddict path")
}
