db = db.getSiblingDB('admin');
db.createUser({
  user: "derrickk",
  mechanisms: ['MONGODB-X509'],
  credentials: {
    certificate: base64EncodeCertFile('ad.pub')
  },
  roles: [{role: "dbAdminAnyDatabase", db: "admin"}, { role: 'readWriteAnyDatabase', db: "admin" }]
});

if (!db.getCollectionNames().includes('servers')) {
  db.createCollection('servers');
}

db = db.getSiblingDB('servers');
db.createUser({
  user: 'microservice_usr',
  mechanisms: ['MONGODB-X509'],
  credentials: {
    certificate: base64EncodeCertFile('micro_usr.pub')
  },
  roles: [{ role: 'readWrite', db: 'servers' }]
});

db.shutdownServer();

function base64EncodeCertFile(certFile) {
  const cert = cat(certFile);
  return new BinData(0, cert.toString('base64'));
}
