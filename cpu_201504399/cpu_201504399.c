#include <linux/proc_fs.h>
#include <linux/kernel.h>
#include <linux/module.h>
#include <linux/init.h>
#include <linux/sched/signal.h>
#include <linux/seq_file.h>
#include <linux/sched.h>
#include <linux/fs.h>

struct task_struct *task; //estructura definifa en sched.h para tareas /procesos 
struct task_struct *task_child; //estructura necesaria para iterar a través de tareas secundarias
struct list_head *list;  //estructura necesaria para recorrer la lista en cada tarea -> estructura de hijos.

static int escribir_archivo(struct seq_file * archivo, void *v){
    seq_printf(archivo,"***********************************************************\n");
    seq_printf(archivo,"*          201504399         Maria Herrera                *\n");
    seq_printf(archivo,"*          PROYECTO 1 SOPES 1 Vac Dic 2020                *\n");
    seq_printf(archivo,"*                         CPU                             *\n");
    seq_printf(archivo,"***********************************************************\n");
    
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
    proc_create("cpu_201504399",0,NULL,&operaciones);
    printk(KERN_INFO "%s","Cargando modulo.\n");
    printk(KERN_INFO "%s","María de los Angeles Herrera Sumalé.\n");

    for_each_process(task){

        printk(KERN_INFO "\nPADRE PID: %d PROCESO: %s ESTADO: %ld",task->pid, task->comm,task->state);
        list_for_each(list, &task->children){
            task_child = list_entry(list,struct task_struct, sibling);
            printk(KERN_INFO "\n HIJO DE %s[%d] PID: %d PROCESO: %s  ESTADO:  %ld",task->comm, task->pid,task_child->pid,task_child->comm,task_child->state);
            
        }
        printk("*****************************************************");
    }
    return 0;
}

void salir(void){ //modulo salida 

    remove_proc_entry("cpu_201504399",NULL);
    printk(KERN_INFO "%s","Removiendo modulo.\n");
    printk(KERN_INFO "Diciembre 2020.\n");

}

module_init(iniciar);
module_exit(salir);

MODULE_LICENSE("GPL");
MODULE_DESCRIPTION("Itera a través de todos los procesos / procesos hijos en el SO.");
MODULE_AUTHOR("201504399");
