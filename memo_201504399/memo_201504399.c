#include <linux/proc_fs.h>
#include <linux/seq_file.h>
#include <asm/uaccess.h>
#include <linux/hugetlb.h>
#include <linux/module.h>
#include <linux/init.h>
#include <linux/kernel.h>
#include <linux/fs.h>


#include <linux/swap.h>
#include <asm/page.h>
#include <linux/mmzone.h>
#include <linux/cpumask.h>
#include <linux/kernel_stat.h>

#define BUFSIZE  150

MODULE_LICENSE("GPL");
MODULE_DESCRIPTION("Ecribir informacion de la memoria ram.");
MODULE_AUTHOR("201504399");

struct sysinfo inf;

static int escribir_archivo(struct seq_file * archivo, void *v){
    si_meminfo(&inf);
    long total_memoria = (inf.totalram *4 );
    long memoria_libre = (inf.freeram * 4);
    long buffer = (inf.bufferram);
    long cached= (global_node_page_state(NR_FILE_PAGES) * 2 )- inf.bufferram ;
    long memoria_utilizada = total_memoria - (memoria_libre+buffer +cached);
    long porcentaje =((memoria_utilizada *100)/total_memoria);
    //Total Memory - (Free + Buffers + Cached) 

   /* seq_printf(archivo,"***********************************************************\n");
    seq_printf(archivo,"*          201504399         Maria Herrera                *\n");
    seq_printf(archivo,"*          PROYECTO 1 SOPES 1 Vac Dic 2020                *\n");
    seq_printf(archivo,"*                     MEMORIA RAM                         *\n");
    seq_printf(archivo,"***********************************************************\n");
    seq_printf(archivo,"SO Ubuntu 18.04.5\n");
     seq_printf(archivo,"Memoria Libre: \t %8lu KB - %8lu MB\n",memoria_libre, memoria_libre /1024);*/
    seq_printf(archivo,"{\n");
    seq_printf(archivo,"\"Total\": \"%8lu\" \n", total_memoria /1024); //total memoria ram
    seq_printf(archivo,"\"Uso\": \"%8lu\" \n", memoria_utilizada/1024); // total memoria consumida
    seq_printf(archivo,"\"Porcentaje\": \"%li\" \n",porcentaje); //porcentaje de consumo
     seq_printf(archivo,"} \n");
    return 0;
}


static int al_abrir(struct inode *inode, struct file *file){
    return single_open(file,escribir_archivo,NULL);
}

static struct file_operations operaciones = 
{
    .open = al_abrir,
    .read = seq_read
};

int iniciar(void){ //modulo de inicio 
    proc_create("memo_201504399",0,NULL,&operaciones);
    printk(KERN_INFO "%s","Cargando modulo.\n");
    printk(KERN_INFO "%s","201504399\n");
   
    return 0;
}

void salir(void){ //modulo salida 

    remove_proc_entry("memo_201504399",NULL);
    printk(KERN_INFO "%s","Removiendo modulo.\n");
    printk(KERN_INFO "Sistemas Operativos 1.\n");

}

module_init(iniciar);
module_exit(salir);



