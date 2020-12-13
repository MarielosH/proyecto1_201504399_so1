#include <linux/proc_fs.h>
#include <linux/seq_file.h>
#include <asm/uaccess.h>
#include <linux/hugetlb.h>
#include <linux/module.h>
#include <linux/init.h>
#include <linux/kernel.h>
#include <linux/fs.h>

#define BUFSIZE  150

MODULE_LICENSE("GPL");
MODULE_DESCRIPTION("Ecribir informacion de la memoria ram.");
MODULE_AUTHOR("201504399");

struct sysinfo inf;

static int escribir_archivo(struct seq_file * archivo, void *v){
    si_meminfo(&inf);
    long total_memoria = (inf.totalram * 8);
    long memoria_libre = (inf.freeram * 8);
    
    seq_printf(archivo,"***********************************************************\n");
    seq_printf(archivo,"*          201504399         Maria Herrera                *\n");
    seq_printf(archivo,"*          PROYECTO 1 SOPES 1 Vac Dic 2020                *\n");
    seq_printf(archivo,"*                     MEMORIA RAM                         *\n");
    seq_printf(archivo,"***********************************************************\n");
    seq_printf(archivo,"SO Ubuntu 18.04.5\n");
    seq_printf(archivo,"Memoria Total: \t %8lu KB - %8lu MB\n",total_memoria, total_memoria /1024);
    seq_printf(archivo,"Memoria Libre: \t %8lu KB - %8lu MB\n",memoria_libre, memoria_libre /1024);
    seq_printf(archivo,"Memoria en uso: \t %li %%\n",(memoria_libre *100)/total_memoria);
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



